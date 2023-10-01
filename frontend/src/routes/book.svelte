<script lang="ts">
    import { goto } from "$app/navigation";
    export let name: string;
    export let id: string;
    export let cover: string;
    const readBook = () => goto(`${window.origin}/read/${id}`);
</script>

<div on:click={readBook}>
    <img src={cover}>
    <p> {name} </p>
</div>

<style>
    img {
        width: 115px;
        height: 190px;
    }
    div {
        width: 115px;
        cursor: pointer;
        text-align: center;
        height: fit-content;
    }
    p {
        font-size: 17px;
        word-break: break-all;
    }
</style>
