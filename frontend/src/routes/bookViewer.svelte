<script lang="ts">
    import { onMount } from "svelte";
    import { Epub } from "./epub";
    import { callApi } from "./utils";

    function renderBook(bookName: string, div: HTMLElement) {
        callApi(`http://localhost:8080/${bookName}`, "GET", {}).then((json) => {
            let e = new Epub(bookName, json.Files, div);
            e.render();
        });
    }

    onMount(() => {
        renderBook("Dune", document.getElementById("book-view")!);
    });
</script>

<div>
    <p> ebook reader </p>
    <hr>
    <div class="book-view" id="book-view"></div>
</div>

<style>
    #book-view {
        width: 600px;
        margin: 0 auto;
    }
</style>
