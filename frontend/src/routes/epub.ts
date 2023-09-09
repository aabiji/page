import * as utils from "./utils";

export class Epub {
    name: string;
    // array containing urls to the xhtml/html files within the epub
    files: string[];
    // promises for downloading css files required for rendering the epub
    css_promises: Promise<void>[] = [];
    // promises for downloading the epub files
    page_promises: Promise<string>[] = [];
    // HTMLElement used to hold all the rendered epub content
    render_container: HTMLElement;
    // the default css to apply for when the epub's xhtml/html files don't have adequate css
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
        this.render_container = container;
    }

    private correct_document_links(root: HTMLElement, content_type: string) {
        let attrs = {
            img: {"application/xhtml+xml": "image", "text/html": "img"},
            img_src: {"application/xhtml+xml": "xlink:href", "text/html": "src"},
        }

        if (root.nodeName.toLowerCase() == attrs.img[content_type]) {
            let src = root.getAttribute(attrs.img_src[content_type]);
            let url = utils.static_file_url(`${this.name}/${src}`);
            root.setAttribute(attrs.img_src[content_type], url);
        }

        for (let c of root.children) {
            this.correct_document_links(c as HTMLElement, content_type);
        }
    }

    private apply_document_css(meta_tags: HTMLCollection, doc: Document) {
        for (let i = 0; i < meta_tags.length; i++) {
            let node = meta_tags[i];
            let lnode = node as HTMLLinkElement;
            if (node.nodeName.toLowerCase() != "link" || lnode.rel != "stylesheet") continue;

            let css_url = lnode.href.replace(`${window.location.origin}/`, "");
            css_url = utils.static_file_url(`${this.name}/${css_url}`);

            const p = utils.download_file(css_url).then((css: string) => {
                let style = document.createElement("style");
                style.textContent += this.default_css;
                // the default css will be overidden by the epub's css, if there's any
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
            iframe.scrolling = "no";
        }
        this.render_container.appendChild(iframe);
   }

    private render_page(content: string, content_type: string) {
        const fetched_doc = new DOMParser().parseFromString(content, content_type);
        let doc = document.implementation.createHTMLDocument();
        doc.body = fetched_doc.body;
        this.correct_document_links(doc.body, content_type);

        this.apply_document_css(fetched_doc.head.children, doc);
        // Wait until all the css files have been downloaded to process the resulting document.
        Promise.all(this.css_promises).then(() => this.encapsulate_page(doc));
    }

    render() {
        let content_types: string[] = [];
        for (let file of this.files) {
            let url = utils.static_file_url(file);
            const p = utils.download_file(url).then((content: string) => content);
            content_types.push(utils.get_content_type(file));
            this.page_promises.push(p);
        }

        // Wait until all the html files have been downloaded to render the pages in order
        Promise.all(this.page_promises).then((html_pages) => {
            for (let i = 0; i < html_pages.length; i++) {
                this.render_page(html_pages[i], content_types[i]);
            }
        });
    }
}
