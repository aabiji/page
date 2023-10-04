<script lang="ts">
    import { goto } from "$app/navigation";
    import * as utils from "$lib/utils";

    export let id: number;
    export let title: string;
    export let cover: string;
    const readBook = () => goto(`${window.origin}/read/${id}`);
</script>

<div class="container">
    <div on:keypress={readBook} tabindex=0 role="button" on:click={readBook}>
        <img alt="Cover" src={cover}>
    </div>
    <p> {title} </p>
</div>

<style>
    img {
        width: 115px;
        height: 190px;
        cursor: pointer;
    }

    p {
        font-size: 17px;
        word-break: break-all;
    }

    .container {
        width: 125px;
        text-align: center;
        height: fit-content;
    }
</style>