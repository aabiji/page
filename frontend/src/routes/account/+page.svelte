<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import * as utils from "$lib/utils";
    import Navbar from "../../components/navbar.svelte";

    function deleteAccount() {
        let url = `${utils.backendOrigin}/user/delete`;
        utils.callApi(url, "POST").then((response) => {
            if (utils.serverError in response) return;
            utils.removeCookie("userId");
            localStorage.clear();
            goto("/auth");
        });
    }

    onMount(() => {
        utils.redirectIfNotAuth();
    });
</script>

<Navbar />
<div class="container">
    <h1> App settings </h1>
    <hr>
    <h3> Account </h3>
    <button on:click={deleteAccount}> Delete account </button>
</div>

<style>
    hr {
        height: 1px;
        border: 1px solid var(--background-accent);
    }

    button {
        color: white;
        font-size: 18px;
        padding: 5px 5px;
        background-color: red;
    }
   
    .container {
        padding: 10px;
        margin-top: var(--navbar-height);
    }
</style>