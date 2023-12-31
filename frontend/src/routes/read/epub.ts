import * as utils from "$lib/utils";

export class EpubViewer {
    // HTMl/XHTML files within the epub file.
    files: string[];
    // Vertical scrolling offsets within the HTML/XHTML files.
    scrolls: number[];
    // The amount to vertically scroll at once.
    // Is set to the height of the renderContainer minus some padding.
    scrollStep: number;
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

    constructor(scrollOffsets: number[], files: string[], pageIdx: number, container?: HTMLElement) {
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

    private getImagesInDocument(doc: Document): HTMLElement[] {
        let images = Array.from(doc.getElementsByTagName("image"));
        let imgs = Array.from(doc.getElementsByTagName("img"));

        let all_imgs = [];
        all_imgs.push(...images);
        all_imgs.push(...imgs);
        return all_imgs as HTMLElement[];
    }

    private injectDefaultCSS(doc: Document) {
        let style = doc.createElement("style");
        style.textContent = this.defaultCss;
        doc.head.appendChild(style);
    }

    jumpToSection(sectionPath: string) {
        let sectionParts = sectionPath.split("#");
        let section = sectionParts[1] == undefined ? "" : sectionParts[1];
        let index = this.files.indexOf(sectionParts[0]);
        if (index == -1) return;
        this.scrolls[this.pageIdx] = 0;
        this.pageIdx = index;
        this.render(section);
    }

    private correctLinks(doc: Document) {
        let sourceAttr = doc.getElementsByTagName("img").length > 0 ? "src" : "xlink:href";

        let images = this.getImagesInDocument(doc);
        for (let image of images) {
            let source = image.getAttribute(sourceAttr)!;
            image.setAttribute(sourceAttr, utils.staticFileUrl(source));
        }

        // NOTE: for now, all links are disabled
        let links = doc.getElementsByTagName("a");
        for (let link of links) link.removeAttribute("href");
    }

    private getLastImageWithinRange(iframe: HTMLIFrameElement, end: number): HTMLImageElement {
        let images = this.getImagesInDocument(iframe.contentWindow!.document);
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

        let parentTags = ["svg", "div", "picture"];
        let parent = lastImage.parentNode as HTMLElement;
        let nodeToResize = parent != null && parentTags.includes(parent.nodeName) ? parent : lastImage;

        nodeToResize.style.height = `${adjustedHeight}px`;
        nodeToResize.style.width = "auto";
    }

    private scrollCurrentPage(iframe: HTMLIFrameElement, event: MouseEvent) {
        let scrollOffset = this.scrolls[this.pageIdx];
        const scrollDirection = event.clientX > this.containerMidPoint ? 1 : -1;
        scrollOffset += this.scrollStep * scrollDirection;
        scrollOffset = Math.max(-1, Math.min(scrollOffset, this.docHeight(iframe)));
        iframe.contentWindow!.scrollTo({top: scrollOffset, behavior: "smooth"});
        this.scrolls[this.pageIdx] = scrollOffset;

        const overflow = scrollOffset < 0 || scrollOffset >= this.docHeight(iframe);
        if (!overflow) {
            this.adjustLastImageHeight(iframe);
            return; // No need to change the current page
        }

        const pageDirection = scrollOffset >= this.docHeight(iframe) ? 1 : -1;
        this.pageIdx = Math.max(0, Math.min(this.pageIdx + pageDirection, this.files.length-1));
        this.render();
    }

    private renderPage(content: string, elementId: string) {
        let contentType = this.files[this.pageIdx].endsWith(".html") ? "text/html" : "application/xhtml+xml";
        let doc = new DOMParser().parseFromString(content, contentType as DOMParserSupportedType);
        this.injectDefaultCSS(doc);
        this.correctLinks(doc);

        let iframe = document.createElement("iframe");
        iframe.srcdoc = doc.documentElement.innerHTML;
        iframe.scrolling = "no";
        iframe.onload = () => {
            iframe.style.height = "inherit";
            iframe.contentDocument!.addEventListener("click", (event) => {
                this.scrollCurrentPage(iframe, event);
            });
            this.adjustLastImageHeight(iframe);

            let iwindow = iframe.contentWindow!;
            let targetElement = iwindow.document.getElementById(elementId);
            if (targetElement != undefined) {
                this.scrolls[this.pageIdx] = targetElement.getBoundingClientRect().y;
            }
            iwindow.scrollTo({top: this.scrolls[this.pageIdx], behavior: "smooth"});
        }

        return iframe;
    }

    // elementId specifies the optional element to jump to when rendering iframe contents.
    async render(elementId: string = "") {
        this.renderContainer.innerHTML = "";
        let url = utils.staticFileUrl(this.files[this.pageIdx]);
        const html = await utils.downloadFile(url);
        let view = this.renderPage(html, elementId);
        this.renderContainer.appendChild(view);
    }
}
