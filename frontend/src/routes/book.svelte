<script lang="ts">
    import { goto } from "$app/navigation";
    export let name: string;
    export let userBookId: string;
    const readBook = () => goto(`${window.origin}/read/${userBookId}`);
</script>

<div on:click={readBook}>
    <img src="default-cover-image.png">
    <p> {name} </p>
</div>

<style>
    img {
        width: 100px;
        height: 175px;
    }
    div {
        width: 100px;
        cursor: pointer;
        text-align: center;
        height: fit-content;
        background-color: red;
    }
    p {
        font-size: 17px;
        margin-top: -1px;
    }
</style>
