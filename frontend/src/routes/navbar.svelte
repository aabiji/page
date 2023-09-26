<script lang="ts">
    let userProfilePic = "default-cover-image.png";
</script>

<div class="container">
    <div class="bar">
        <h1 title="Home" class="logo"><a href="/">Page</a></h1> 
        <input title="Search" type="text" placeholder="Search books">
        <a title="Account" href="/account"><img src={userProfilePic} /></a>
    </div>
    <hr>
</div>

<style>
    .container {
        width: 100%;
        background-color: var(--background-accent);
    }
    .bar * {
        display: inline;
    }
    .bar {
        width: 100%;
        height: 50px;
        text-align: center;
        margin-top: -5px;
    }
    h1 {
        float: left;
        margin-top: 1px;
        margin-left: 10px;
        margin-bottom: 0px;
    }
    img {
        width: 35px;
        height: 35px;
        border-radius: 50%;
        float: right;
        margin-top: 5px;
        margin-right: 20px;
    }
    input {
        width: 50%;
        color: white;
        margin-top: 3px;
        font-size: 18px;
        padding: 10px 10px;
        margin-bottom: 15px;
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
