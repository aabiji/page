<script lang="ts">
    import { onMount } from "svelte";
    import * as utils from "$lib/utils";
    import Navbar from "./navbar.svelte";
    import Book from "./book.svelte";

    let bookIds: number[] = [];
    let fileInput: HTMLElement;
    function uploadFile(event: any) {
        const file = event.target.files[0];
        const formData = new FormData();
        formData.append("file", file);
        let url = "http://localhost:8080/user/book/upload";
        utils.callApi(url, "POST", formData, true).then((response) => {
            utils.cacheBookId(response.BookId);
            bookIds = utils.getBookIdCache();
        });
    }

    onMount(() => {
        utils.redirectIfNotAuth();
        bookIds = utils.getBookIdCache();
    });
</script>

<Navbar />
<div class="container">
    <div class="top">
        <h1> Your books </h1>
        <input on:change={uploadFile} bind:this={fileInput} type="file" style="display: none;">
        <button on:click={() => fileInput.click()}> Upload book </button>
    </div>
    <div class="collection">
        {#each bookIds as bookId}
            <Book name={bookId} id={bookId} />
        {/each}
    </div>
</div>

<style>
    button {
        color: white;
        margin-left: 10px;
        font-size: 18px;
        padding: 5px 5px;
        margin-bottom: 15px;
        background-color: var(--accent-color);
    }
    button:hover {
        background-color: var(--accent-color-darken);
    }
    .container {
        margin-top: 65px;
    }
    .top * {
        display: inline;
    }
    .top {
        margin-bottom: 20px;
    }
    .collection {
        gap: 30px;
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(115px, 1fr));
    }
</style>
