import * as utils from "./utils";

export class Epub {
    name: string;
    url: string;
    files: string[];
    is_html: boolean;
    container: HTMLElement;
    content_type: DOMParserSupportedType;
    attributes: object; // attribute names for the different content types

    constructor(name: string, base_url: string, files: string[], container_div: HTMLElement) {
        this.name = name;
        this.files = files;
        this.url = base_url;
        this.container = container_div;
        this.content_type = utils.get_content_type(this.files[0]);
        this.is_html = this.content_type == "text/html";
        this.attributes = {
            "text/html": {"href": "href"},
            "application/xhtml+xml": {"href": "xlink:href"}
        }
    }

    private get_attr(id: string): string {
        return this.attributes[this.content_type as keyof typeof this.attributes][id];
    }

    private process_node(node: Element): HTMLElement | undefined {
        if (node.innerHTML.split("=").length >= 6) { // todo: this seems somewhat sketchy
            return document.createElement("hr");
        }

        switch (node.nodeName) {
        case "image":
            let attr = this.get_attr("href");
            console.log(node.getAttribute(attr));
            break;

        case "p":
            let p = document.createElement("p");
            p.classList.add("book-text");
            p.appendChild(document.createTextNode(node.innerHTML));
            return p;
        }

        return undefined; // Unrecognized node
    }

    private render_node(node: Element) {
        let n = this.process_node(node);
        if (n != undefined)
            this.container.appendChild(n);

        for (let n of node.children) {
            this.render_node(n);
        }
    }

    // FIXME: the reason why files aren't rendered sequentially
    // is because they all take different times to be fetched.
    render_book() {
        // Render each xhtml/html page:
        for (let file of this.files) {
            let file_url = utils.static_file_url(file);
            utils.call_api(file_url).then((content: string) => {
                const doc = new DOMParser().parseFromString(content, this.content_type);
                this.render_node(doc.body);
            });
        }
    }
}
