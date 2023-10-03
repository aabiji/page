<script lang="ts">
    import { onMount } from "svelte";
    import { writable } from "svelte/store";
 
    import Book from "../components/book.svelte";
    import Navbar from "../components/navbar.svelte";
    import Upload from "../components/upload.svelte";
    import * as utils from "$lib/utils";
 
    interface BookDisplayInfo {
        id: number,
        cover: string,
        title: string,
    };
 
    let books = writable<BookDisplayInfo[]>([]);
    function loadBooks(bookIds: number[]) {
        $books = [];
        for (let i = 0; i < bookIds.length; i++) {
            let id = bookIds[i];
            let obj: utils.Book = utils.cacheGet(utils.BookKey(id));
            let display: BookDisplayInfo = {
                id: id,
                title: obj.Info.Title,
                cover: obj.CoverImagePath,    
            };
            $books.push(display);
        }
    }

    onMount(() => {
        utils.redirectIfNotAuth();
        loadBooks(utils.cacheGet(utils.BooksKey));
    });
</script>

<Navbar />
<div class="container">
    <h1> Your books </h1>
    {#if $books.length == 0}
        <p> Looks like you don't have any books yet! </p>
    {/if}
    <div class="collection">
        {#each $books as b}
            <Book cover={b.cover} title={b.title} id={b.id} />
        {/each}
    </div>
</div>
<Upload />

<style>
    .container {
        padding: 10px;
        margin-top: var(--navbar-height);
    }
    .collection {
        gap: 30px;
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(115px, 1fr));
    }
</style>