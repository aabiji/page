import * as utils from "./utils";

interface File {
    Path: string;
}

export class EpubViewer {
    // HTMl/XHTML files within the epub file.
    files: File[];
    // Vertical scrolling offsets within the HTML/XHTML files.
    scrolls: number[];
    // The amount to vertically scroll at once.
    // Is set to the height of the renderContainer minus some padding.
    scrollStep: number;
    sections: string[] = [];
    // The amount of padding to remove from scrollStep
    // to ensure that text near the document edges remain fully visible.
    pad: number;
    // Index to the currently rendered HTML/XHTML file.
    pageIdx: number;
    // The horizantal middle of the renderContainer in pixels.
    containerMidPoint: number;
    // HTMLElement used to hold all the rendered epub content.
    renderContainer: HTMLElement;
    // The default CSS to apply for when the epub's XHTML/HTML files don't have adequate CSS.
    defaultCss: string = `
        body {
            line-height: 2.0;
            text-indent: 25px;
            text-align: left;
            background-color: #1c1c1c;
            font-family: Arial, Helvetica, sans-serif;
        }
        body :not(a) {
            color: white !important;
            }
        a {
            color: #4287f5 !important;
        }
        a:hover {
            color: #5a98fa !important;
        }
        pre, code {
            color: white !important;
            background-color: #101010 !important;
        }
    `;

    constructor(scrollOffsets: number[], files: File[], pageIdx: number, container?: HTMLElement) {
        this.pad = 10;
        this.files = files;
        this.pageIdx = pageIdx;
        this.scrolls = scrollOffsets;
        this.renderContainer = container!;
        if (this.renderContainer != undefined) {
            this.scrollStep = this.renderContainer.clientHeight - this.pad;
            this.containerMidPoint = this.renderContainer.clientWidth / 2;
        } else {
            this.scrollStep = 0;
            this.containerMidPoint = 0;
        }
    }

    private docHeight(iframe: HTMLIFrameElement) {
        return iframe.contentWindow!.document.documentElement.scrollHeight - this.pad;
    }

    private imageAttributes() {
        return this.files[this.pageIdx].Path.endsWith(".html")
                 ? {tag: "img", source: "src"}
                 : {tag: "image", source: "xlink:href"};
    }

    private injectDefaultCSS(doc: Document) {
        let style = document.createElement("style");
        style.textContent = this.defaultCss;
        doc.head.appendChild(style);
    }

    private correctImageLinks(doc: Document) {
        let attr = this.imageAttributes();
        let images = doc.getElementsByTagName(attr.tag);
        for (let image of images) {
            let source = image.getAttribute(attr.source)!;
            image.setAttribute(attr.source, utils.staticFileUrl(source));
        }
    }

    private getLastImageWithinRange(iframe: HTMLIFrameElement, end: number): HTMLImageElement {
        let attr = this.imageAttributes();
        let images = iframe.contentWindow!.document.getElementsByTagName(attr.tag);
        let imagesWithinRange = Array.from(images).filter((img) => {
            return img.getBoundingClientRect().top < end;
        });
        let lastImage = imagesWithinRange[imagesWithinRange.length - 1];
        return lastImage as HTMLImageElement;;
    }

    // Adjusting the image height to ensure that it remains fully
    // visible within the scrollable area, preventing any clipping during scrolling.
    private adjustLastImageHeight(iframe: HTMLIFrameElement) {
        let end = this.scrolls[this.pageIdx] + this.scrollStep;
        if (end > this.docHeight(iframe)) end = this.scrollStep;

        let lastImage = this.getLastImageWithinRange(iframe, end);
        if (lastImage == undefined) return;

        let rect = lastImage.getBoundingClientRect();
        let overflow = Math.max(0, rect.bottom - end);
        let adjustedHeight = rect.height - overflow;

        let parent = lastImage.parentNode as HTMLElement;
        let nodeToResize = parent != null && parent.nodeName == "svg" ? parent : lastImage;
        nodeToResize.style.height = `${adjustedHeight}px`;
        nodeToResize.style.width = "auto";
    }

    private scrollCurrentPage(iframe: HTMLIFrameElement, event: MouseEvent) {
        let scrollOffset = this.scrolls[this.pageIdx];
        const scrollDirection = event.clientX > this.containerMidPoint ? 1 : -1;
        scrollOffset += this.scrollStep * scrollDirection;
        scrollOffset = Math.max(-1, Math.min(scrollOffset, this.docHeight(iframe)));
        iframe.contentWindow!.document.documentElement.scrollTo(0, scrollOffset);
        this.sections[this.pageIdx] = "";
        this.scrolls[this.pageIdx] = scrollOffset;

        const overflow = scrollOffset < 0 || scrollOffset >= this.docHeight(iframe);
        if (!overflow) {
            this.adjustLastImageHeight(iframe);
            return; // No need to change the current page
        }

        const pageDirection = scrollOffset >= this.docHeight(iframe) ? 1 : -1;
        this.pageIdx = Math.max(0, Math.min(this.pageIdx + pageDirection, this.files.length - 1));
        this.render();
    }

    private scrollTo(targetWindow: Window, y: number, duration: number) {
        const startY = targetWindow.scrollY;
        const startTime = performance.now();
        const step = (timestamp: number) => {
            const currentTime = timestamp - startTime;
            const progress = Math.min(currentTime / duration, 1);
            const newY = startY + (y - startY) * progress;

            if (progress < 0) return;
            targetWindow.scrollTo(0, newY);
            if (progress < 1) targetWindow.requestAnimationFrame(step);
        }
        targetWindow.requestAnimationFrame(step);
    }

    private renderPage(content: string) {
        let contentType = this.files[this.pageIdx].Path.endsWith(".html") ? "text/html" : "application/xhtml+xml";
        const doc = new DOMParser().parseFromString(content, contentType as DOMParserSupportedType);
        this.injectDefaultCSS(doc);
        this.correctImageLinks(doc);

        let iframe = document.createElement("iframe");
        iframe.srcdoc = doc.documentElement.innerHTML;
        iframe.scrolling = "no";
        iframe.onload = () => {
            iframe.style.height = "inherit";
            iframe.contentDocument!.addEventListener("click", (event) => {
                this.scrollCurrentPage(iframe, event);
            });

            this.adjustLastImageHeight(iframe);
            let doc = iframe.contentWindow!.document;
            this.scrollTo(iframe.contentWindow!, this.scrolls[this.pageIdx], 1000);

            let section = doc.getElementById(this.sections[this.pageIdx]);
            if (section != undefined) section.scrollIntoView({behavior: "smooth"});
        }

        return iframe;
    }

    async render() {
        this.renderContainer.innerHTML = "";
        let url = utils.staticFileUrl(this.files[this.pageIdx].Path);
        const html = await utils.downloadFile(url);
        let view = this.renderPage(html);
        this.renderContainer.appendChild(view);
    }
}
