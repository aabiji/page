import * as utils from "./utils";

export class Epub {
    name: string;
    files: string[];
    container: HTMLElement;
    content_type: DOMParserSupportedType;
    attributes: object; // attribute names for the different content types

    css_promises: Promise<void>[];
    page_promises: Promise<string>[];

    constructor(name: string, files: string[], container_div: HTMLElement) {
        this.css_promises = [];
        this.page_promises = [];

        this.name = name;
        this.files = files;
        this.container = container_div;
        this.content_type = utils.get_content_type(this.files[0]);
        this.attributes = {
            "text/html": {"href": "href", "img": "img"},
            "application/xhtml+xml": {"href": "xlink:href", "img": "image"}
        }
    }

    private get_attr(id: string): string {
        return this.attributes[this.content_type as keyof typeof this.attributes][id];
    }

    private parse_node(node: Element): HTMLElement | undefined {
        for (let c of Array.from(node.classList)) {
            if (c.includes("calibre") && node.innerHTML.split("=").length >= 6) {
                return document.createElement("hr");
            }
        }

        switch (node.nodeName.toLowerCase()) {
        case this.get_attr("img"):
            let href = node.getAttribute(this.get_attr("href"));
            let real_href = utils.static_file_url(`${this.name}/${href}`);
            let img = document.createElement("img");
            img.src = real_href;
            return img;

        case "a":
            let a = document.createElement("a");
            a.innerHTML = node.innerHTML;
            a.href = node.getAttribute(this.get_attr("href"))!;
            return a;

        case "p":
        case "h1": case "h2": case "h3": case "h4": case "h5": case "h6":
            let text = document.createElement(node.nodeName);
            if (node.children.length == 0) {
                text.appendChild(document.createTextNode(node.innerHTML));
            }
            return text;

        default:
            return undefined; // TODO: Unrecognized node: div ...
        }
    }

    private process_node(node: Element, doc: HTMLElement) {
        let n = this.parse_node(node);
        if (n != undefined)
            doc.appendChild(n);

        for (let n of node.children) {
            this.process_node(n, doc);
        }
    }

    // we're combining all the css into one file
    // maybe we could read the css of one page and extrapolate that to the rest
    // of the book ??? (for efficiency reasons)
    private apply_page_css(meta_tags: HTMLCollection, doc: Document) {
        for (let i = 0; i < meta_tags.length; i++) {
            let node = meta_tags[i];
            if (node.nodeName.toLowerCase() != "link") continue;
            if ((node as HTMLLinkElement).rel != "stylesheet") continue;

            let css_url = (node as HTMLLinkElement).href.replace(`${window.location.origin}/`, "");
            css_url = utils.static_file_url(`${this.name}/${css_url}`);

            const p = utils.download_file(css_url).then((css: string) => {
                let style = document.createElement("style");
                style.textContent += css;
                doc.head.appendChild(style);
            });
            this.css_promises.push(p);
        }
    }

    private encapsulate_page(doc: Document) {
        let iframe = document.createElement("iframe");
        iframe.srcdoc = doc.documentElement.innerHTML;
        iframe.onload = () => { // Resize iframe height to fit content
            let h = iframe.contentWindow!.document.documentElement.scrollHeight;
            iframe.style.height = `${h}px`;
        }
        this.container.appendChild(iframe);
   }

    private render_page(content: string) {
        const parsed_doc = new DOMParser().parseFromString(content, this.content_type);

        let doc = document.implementation.createHTMLDocument();
        this.process_node(parsed_doc.body, doc.body);

        this.apply_page_css(parsed_doc.head.children, doc);
        // Wait until all the css files have been downloaded to process the resulting document.
        Promise.all(this.css_promises).then(() => {
            this.encapsulate_page(doc)
            this.css_promises = [];
        });
    }

    render() {
        for (let i = 0; i < this.files.length; i++) {
            const file = this.files[i];
            let url = utils.static_file_url(file);
            const p = utils.download_file(url).then((content: string) => content );
            this.page_promises.push(p);
        }

        // Wait until all the html files have been downloaded to render the pages in order
        Promise.all(this.page_promises).then((html_pages) => {
            for (let html of html_pages) {
                this.render_page(html);
            }
            this.page_promises = [];
        });
    }
}
