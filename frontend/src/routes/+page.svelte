<script lang="ts">
    import { onMount } from 'svelte';

    async function call_api(url: string) {
        const response = await fetch(url, {
            method: "GET",
            mode: "cors",
        });
        return await response.text();
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
                let text = node.innerHTML.replace(/\-\?/, "-")
                p.appendChild(document.createTextNode(text));
                container.appendChild(p)
                break;
            }
        }
    }

    onMount(() => {
        const url = "http://localhost:8080/static/Dune/OEBPS/part2_split_000.xhtml";

        call_api(url).then((content) => {
            let doc = new DOMParser().parseFromString(content, "text/xml");
            render_xhtml(document.getElementById("book-view")!, doc);
        });
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
    :global(p) {
        line-height: 2.0;
        text-indent: 50px;
    }
</style>
