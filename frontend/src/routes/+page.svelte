<script lang="ts">
    import { onMount } from 'svelte';

    async function call_api(url: string) {
        const response = await fetch(url, {
            method: "GET",
            mode: "cors",
        });
        return await response.text();
    }

    function get_content_type(file_url: string): DOMParserSupportedType {
        // all content types:
        // ["application/xhtml+xml","application/xml","image/svg+xml","text/html","text/xml"]
        let extention = file_url.split(".")[1];
        return extention == "xhtml" ? "application/xhtml+xml" : "text/html";
    }

    function render_xhtml(container: HTMLElement, doc: Document) {
        for (let node of doc.body.children) {
            if (node.innerHTML.split("=").length >= 6) {
                container.appendChild(document.createElement("hr"));
                continue;
            }

            switch (node.nodeName) {
            case "p":
                let p = document.createElement("p");
                p.classList.add("book-text");
                let text = node.innerHTML.replace(/\-\?/, "-")
                p.appendChild(document.createTextNode(text));
                container.appendChild(p)
                break;
            }
        }
    }

    function render_book(book_name: string) {
        const url = `http://localhost:8080/${book_name}`;
        call_api(url).then((content) => {
            let view = document.getElementById("book-view")!;

            let c = content.slice(1, content.length - 1); // remove '[' and ']'
            let filepaths = c.split(" ");

            for (let file of filepaths) {
                let netpath = file.replace("BOOKS", "static");
                let file_url = `http://localhost:8080/${netpath}`;

                call_api(file_url).then((content) => {
                    let ct = get_content_type(file_url);
                    let doc = new DOMParser().parseFromString(content, ct);
                    render_xhtml(view, doc);
                });
            }
        });
    }

    onMount(() => {
        render_book("Dune");
    });
</script>

<div>
    <h1> ebook reader </h1>
    <hr>
    <div id="book-view"></div>
</div>

<style>
    #book-view {
        width: 600px;
        margin: 0 auto;
    }
</style>
