async function refreshApiKey() {
    const url = "https://api.youxiaohou.com/config/v2?ver=6.2.7&a=YouXiaoHou";
    const response = await fetch(url, { method: "POST" });
    const tx = Buffer.from(await response.text(), "base64").toString("utf-8");
    const config = JSON.parse(tx);
    const Baidu_Verify_url = (config.pcs)?.["3"];
    if (!Baidu_Verify_url) {
        throw new Error("Baidu_Verify_url not found in config");
    } else {
        await Bun.write(new URL('./baidu_verify_url.txt', import.meta.url), Baidu_Verify_url);
        console.log('Baidu_Verify_url refreshed and saved to baidu_verify_url.txt:', Baidu_Verify_url);
        return Baidu_Verify_url;
    }
}

await refreshApiKey();