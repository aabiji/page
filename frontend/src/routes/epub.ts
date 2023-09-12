import * as utils from "./utils";

export class EpubViewer {
    files: object[];
    pageIndex: number;
    // HTMLElement used to hold all the rendered epub content.
    renderContainer: HTMLElement;
    // The default css to apply for when the epub's xhtml/html files don't have adequate css.
    defaultCss: string = `
        body {
            color: black;
            line-height: 2.0;
            text-indent: 25px;
            text-align: left;
            background-color: white;
        }
    `;

    constructor(files: object[], container: HTMLElement) {
        this.files = files;
        this.pageIndex = 0;
        this.renderContainer = container;
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

    private renderPage(content: string, contentType: string) {
        const doc = new DOMParser().parseFromString(content, contentType as DOMParserSupportedType);
        this.injectDefaultCSS(doc);
        this.correctImageLinks(doc);

        let iframe = document.createElement("iframe");
        iframe.scrolling = "no";
        iframe.srcdoc = doc.documentElement.innerHTML;
        iframe.onload = () => { // Resize iframe height to fit content
            iframe.style.height = "inherit";

            let file = this.files[this.pageIndex];
            const half = this.renderContainer.clientWidth / 2;
            const scrollStep = this.renderContainer.clientHeight;
            let docHeight = iframe.contentWindow!.document.documentElement.scrollHeight;
            iframe.contentDocument!.addEventListener("click", (event) => {
                let scrollDirection = event.clientX > half ? 1 : -1;
                file.ScrollOffset += (scrollStep * scrollDirection);
                iframe.contentWindow.document.documentElement.scrollTo(0, file.ScrollOffset);

                if (file.ScrollOffset >= 0 && file.ScrollOffset < docHeight) return;
                let fileCount = this.files.length;
                let pageDirection = file.ScrollOffset >= docHeight ? 1 : -1;
                this.pageIndex += pageDirection;
                if (this.pageIndex < 0) this.pageIndex = 0;
                if (this.pageIndex >= fileCount) this.pageIndex = fileCount;
                this.renderContainer.innerHTML = "";
                this.render();
            });
        }

        return iframe;
    }

    render() {
        let file = this.files[this.pageIndex];
        let url = utils.staticFileUrl(file.Path);
        utils.downloadFile(url).then((content: string) => {
            let view = this.renderPage(content, file.ContentType);
            this.renderContainer.appendChild(view);
        });
    }
}
