import * as utils from "./utils";

export class Epub {
    name: string;
    url: string;
    files: string[];
    container: HTMLElement;
    content_type: DOMParserSupportedType;
    attributes: object; // attribute names for the different content types

    constructor(name: string, base_url: string, files: string[], container_div: HTMLElement) {
        this.name = name;
        this.files = files;
        this.url = base_url;
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

    private process_node(node: Element): HTMLElement | undefined {
        for (let c of Array.from(node.classList)) {
            if (c.includes("calibre") && node.innerHTML.split("=").length >= 6) {
                return document.createElement("hr");
            }
        }

        switch (node.nodeName) {
        case this.get_attr("img"):
            let href = node.getAttribute(this.get_attr("href"));
            let real_href = utils.static_file_url(`${this.name}/${href}`);
            let img = document.createElement("img");
            img.src = real_href;
            return img;

        case "p":
        case "h1": case "h2": case "h3": case "h4": case "h5": case "h6":
            let text = document.createElement(node.nodeName);
            text.appendChild(document.createTextNode(node.innerHTML));
            return text;
        }

        return undefined; // TODO: Unrecognized node: div ...
    }

    private render_node(node: Element) {
        let n = this.process_node(node);
        if (n != undefined)
            this.container.appendChild(n);

        for (let n of node.children) {
            this.render_node(n);
        }
    }

    render() {
        // FIXME: we have this problem where the files are being rendered out of
        // order. This is because the files don't finish downloading in
        // chronological order ...
        // Render each xhtml/html page:
        for (let file of this.files) {
            let url = utils.static_file_url(file);
            utils.download_file(url).then((content: string) => {
                const doc = new DOMParser().parseFromString(content, this.content_type);
                this.render_node(doc.body);
            });
        }
    }
}
