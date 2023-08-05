const {args:[url]} = Deno;

function isUrlValid(input: string) {
    const regex = new RegExp("^(https?://)?(www\\.)?([-a-z0-9]{1,63}\\.)*?[a-z0-9][-a-z0-9]{0,61}[a-z0-9]\\.[a-z]{2,6}(/[-\\w@\\+\\.~#\\?&/=%]*)?$");
    return regex.test(input);
}

async function shorten(url:string):Promise<{result_url:string}>{
    if(url===""||url===undefined){
        throw new Error("URL is empty");
    }
    if(!isUrlValid(url)){
        throw new Error("URL is not valid");
    }

    const options = {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({url}),
    };
    const res=await fetch("https://cleanuri.com/api/v1/shorten",options);
    const json=await res.json();
    return json;
}

try{
    const result=await shorten(url);
    console.log("Shortened URL:");
    console.log(result.result_url);

    await Deno.writeTextFile("shortened.txt",`${result.result_url}->${url}\n`,{append:true});
    console.log("Shortened URL saved to shortened.txt");
}catch (e) {
    console.log(e.message);
}