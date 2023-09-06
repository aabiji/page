<script lang="ts">
    import { onMount } from "svelte";
    import { Epub } from "./epub";

    onMount(() => {
        let c = document.getElementById("book-view");
        // actually get files:
        //let c = content.slice(1, content.length - 1); // remove '[' and ']'
        //let filepaths = c.split(" ");
        let e = new Epub("Dune","http://localhost:8080/Dune",["BOOKS/Dune/titlepage.xhtml"],c);
        e.render_book();
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
