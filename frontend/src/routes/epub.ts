import * as utils from "./utils";

export class Epub {
    name: string;
    // Array containing all renderable files within the epub. ex.
    // [{Path: "", ContentType: ""}, ...]
    files: string[];
    // Promises for downloading the epub files.
    pagePromises: Promise<string>[] = [];
    // HTMLElement used to hold all the rendered epub content.
    renderContainer: HTMLElement;
    // The default css to apply for when the epub's xhtml/html files don't have adequate css.
    default_css: string = `
        body {
            color: black;
            line-height: 2.0;
            text-indent: 25px;
            text-align: left;
            background-color: white;
        }
    `;

    constructor(name: string, files: string[], container: HTMLElement) {
        this.name = name;
        this.files = files;
        this.renderContainer = container;
    }

    private renderPage(content: string, contentType: string) {
        const doc = new DOMParser().parseFromString(content, contentType as DOMParserSupportedType);
        let style = document.createElement("style");
        style.textContent = this.default_css;
        doc.head.appendChild(style);

        let imgSrc = contentType == "text/html" ? "src" : "xlink:href";
        let imgTag = contentType == "text/html" ? "img" : "image";
        let imgs = doc.getElementsByTagName(imgTag);
        for (let img of imgs) {
            let src = img.getAttribute(imgSrc)!;
            img.setAttribute(imgSrc, utils.staticFileUrl(src));
        }

        let iframe = document.createElement("iframe");
        iframe.srcdoc = doc.documentElement.innerHTML;
        iframe.onload = () => { // Resize iframe height to fit content
            let h = iframe.contentWindow!.document.documentElement.scrollHeight;
            iframe.style.height = `${h}px`;
            iframe.scrolling = "no";
        }
        this.renderContainer.appendChild(iframe);
    }

    render() {
        for (let file of this.files) {
            let url = utils.staticFileUrl(file.Path);
            const p = utils.downloadFile(url).then((content: string) => content);
            this.pagePromises.push(p);
        }

        // Wait until all the html files have been downloaded to render the pages in order
        Promise.all(this.pagePromises).then((html_pages) => {
            for (let i = 0; i < html_pages.length; i++) {
                this.renderPage(html_pages[i], this.files[i].ContentType);
            }
        });
    }
}
