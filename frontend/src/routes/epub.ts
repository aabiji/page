import * as utils from "./utils";

export class EpubViewer {
    // HTMl/XHTML files within the epub file.
    files: object[];
    // Vertical scrolling offsets within the HTML/XHTML files.
    scrollOffsets: number[];
    // The amount to vertically scroll at once.
    // Is set to the height of the renderContainer minus some padding.
    scrollStep: number;
    // Index to the currently rendered HTML/XHTML file.
    currentPage: number;
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
        pre, code {
            color: white !important;
            background-color: #101010 !important;
        }
    `;

    constructor(scrollOffsets: number[], files: object[], currentPage: number, container: HTMLElement) {
        this.files = files;
        this.currentPage = currentPage;
        this.scrollOffsets = scrollOffsets;
        this.renderContainer = container;
        this.scrollStep = this.renderContainer.clientHeight - 10;
        this.containerMidPoint = this.renderContainer.clientWidth / 2;
    }

    private correctImageLinks(doc: Document) {
        let imageSource = doc.getElementsByTagName("img").length > 0 ? "src" : "xlink:href";
        let imageTag = imageSource == "src" ? "img" : "image";
        let images = doc.getElementsByTagName(imageTag);
        for (let image of images) {
            let source = image.getAttribute(imageSource)!;
            image.setAttribute(imageSource, utils.staticFileUrl(source));
        }
    }

    private injectDefaultCSS(doc: Document) {
        let style = document.createElement("style");
        style.textContent = this.defaultCss;
        doc.head.appendChild(style);
    }

    private scroll(iframe: HTMLIFrameElement, event: MouseEvent) {
        const docHeight = iframe.contentWindow!.document.documentElement.scrollHeight;
        let scrollOffset = this.scrollOffsets[this.currentPage];

        // Scroll up or down depending on which side the click was registered
        const scrollDirection = event.clientX > this.containerMidPoint ? 1 : -1;
        scrollOffset += this.scrollStep * scrollDirection;

        iframe.contentWindow!.document.documentElement.scrollTo(0, scrollOffset);
        this.scrollOffsets[this.currentPage] = scrollOffset;

        const overflow = scrollOffset < 0 || scrollOffset > docHeight;
        if (!overflow) return; // No need to change current page

        const pageDirection = scrollOffset >= docHeight ? 1 : -1;
        this.currentPage += pageDirection;
        if (this.currentPage < 0)
            this.currentPage = this.files.length - 1;
        else if (this.currentPage == this.files.length)
            this.currentPage = 0;

        this.renderContainer.innerHTML = "";
        this.render();
    }

    private renderPage(content: string, contentType: string) {
        const doc = new DOMParser().parseFromString(content, contentType as DOMParserSupportedType);
        this.injectDefaultCSS(doc);
        this.correctImageLinks(doc);

        let iframe = document.createElement("iframe");
        iframe.srcdoc = doc.documentElement.innerHTML;
        iframe.scrolling = "no";
        iframe.onload = () => {
            iframe.style.height = "inherit";
            iframe.contentDocument!.addEventListener("click", (event) => this.scroll(iframe, event));

            let scrollOffset = this.scrollOffsets[this.currentPage];
            iframe.contentWindow!.document.documentElement.scrollTo(0, scrollOffset);
        }

        return iframe;
    }

    render() {
        let file = this.files[this.currentPage];
        let url = utils.staticFileUrl(file.Path);
        utils.downloadFile(url).then((content: string) => {
            let view = this.renderPage(content, file.ContentType);
            this.renderContainer.appendChild(view);
        });
    }
}
