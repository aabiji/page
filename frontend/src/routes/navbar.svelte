<script lang="ts">
    import { onMount } from "svelte";

    let userProfilePic: string;
    onMount(() => {
        userProfilePic = `${window.origin}/default-profile.png`;
    });
</script>

<div class="bar">
    <h1 title="Home" class="logo"><a href="/">Page</a></h1> 
    <input title="Search" type="text" placeholder="Search books">
    <a title="Account" href="/account"><img alt="Profile photo" src={userProfilePic} /></a>
</div>

<style>
    .bar {
        top: 0;
        height: var(--navbar-height);
        width: 100.5%;
        position: fixed;
        padding: 5px 5px;
        text-align: center;
        margin-left: -10px;
        background-color: var(--background-accent);
    }
    .bar * {
        display: inline;
    }
    h1 {
        float: left;
        margin-top: -2px;
        margin-left: 10px;
        margin-bottom: 0px;
    }
    img {
        width: 35px;
        height: 35px;
        float: right;
        border-radius: 50%;
        margin-right: 20px;
        margin-top: 8px;
    }
    input {
        width: 50%;
        color: white;
        font-size: 18px;
        padding: 8px 8px;
        border: #535454 1px solid;
        background-color: rgba(0,0,0,0);
    }
    input:hover, input:focus {
        outline: none;
        border: var(--accent-color) 1px solid;
    }
    input:-webkit-autofill, input:-webkit-autofill:hover, input:-webkit-autofill:focus {
        font-size: 18px;
        -webkit-text-fill-color: white;
        -webkit-box-shadow: 0 0 0px 40rem #1b1c1c inset;
    }
</style>