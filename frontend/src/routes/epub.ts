import * as utils from "./utils";

export class Epub {
    name: string;
    // array containing urls to the xhtml/html files within the epub
    files: string[];
    // the content type of the files: either application/xhtml+xml,
    // application/xml, image/svg+xml, text/html or text/xml
    content_type: DOMParserSupportedType;
    // promises for downloading css files required for rendering the epub
    css_promises: Promise<void>[] = [];
    // promises for downloading the epub files
    page_promises: Promise<string>[] = [];
    // HTMLElement used to hold all the rendered epub content
    render_container: HTMLElement;

    constructor(name: string, files: string[], container: HTMLElement) {
        this.name = name;
        this.files = files;
        this.render_container = container;
        this.content_type = utils.get_content_type(this.files[0]);
    }

    private get_attr(id: string): string {
        // attribute name mappings for xhtml and html elements
        let attributes: object = {
            "text/html": {"href": "href", "img": "img"},
            "application/xhtml+xml": {"href": "xlink:href", "img": "image"}
        };
        return attributes[this.content_type as keyof typeof attributes][id];
    }

    private process_node(node: Element): HTMLElement | undefined {
        // render horizantal rules for epub files meant to be rendered by the calibre app
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
            return undefined; // unrecognized nodes are undefined -- ex. div, etc.
        }
    }

    private render_node(root: Element, doc: HTMLElement) {
        let n = this.process_node(root);
        if (n != undefined)
            doc.appendChild(n);

        for (let n of root.children) {
            this.render_node(n, doc);
        }
    }

    private apply_rendered_page_css(meta_tags: HTMLCollection, doc: Document) {
        for (let i = 0; i < meta_tags.length; i++) {
            let node = meta_tags[i];
            let lnode = node as HTMLLinkElement;
            if (node.nodeName.toLowerCase() != "link" || lnode.rel != "stylesheet") {
                continue;
            }

            let css_url = lnode.href.replace(`${window.location.origin}/`, "");
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
        this.render_container.appendChild(iframe);
   }

    private render_page(content: string) {
        const parsed_doc = new DOMParser().parseFromString(content, this.content_type);

        let doc = document.implementation.createHTMLDocument();
        this.render_node(parsed_doc.body, doc.body);

        this.apply_rendered_page_css(parsed_doc.head.children, doc);
        // Wait until all the css files have been downloaded to process the resulting document.
        Promise.all(this.css_promises).then(() => this.encapsulate_page(doc));
    }

    render() {
        for (let file of this.files) {
            let url = utils.static_file_url(file);
            const p = utils.download_file(url).then((content: string) => content);
            this.page_promises.push(p);
        }

        // Wait until all the html files have been downloaded to render the pages in order
        Promise.all(this.page_promises).then((html_pages) => {
            for (let html of html_pages) {
                this.render_page(html);
            }
        });
    }
}
