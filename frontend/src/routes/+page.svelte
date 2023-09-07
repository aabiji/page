<script lang="ts">
    import { onMount } from "svelte";

    import { Epub } from "./epub";
    import { call_api } from "./utils";

    function render_book(book_name: string, div: HTMLElement) {
        call_api(`http://localhost:8080/${book_name}`).then((data: string) => {
            let temp = data.slice(1, data.length - 1); // remove '[' and ']'
            let files = temp.split(" ");
            let e = new Epub(book_name, files, div);
            e.render();
        });
    }

    onMount(() => {
        render_book("AnimalFarm", document.getElementById("book-view")!);
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
