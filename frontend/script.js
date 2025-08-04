document.getElementById("sendBtn").addEventListener("click", async () => {
  const userInput = document.getElementById("userInput").value;
  const token = localStorage.getItem("token"); // assume you saved JWT after login

  try {
    const res = await fetch("http://localhost:8080/api/echo", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`
      },
      body: JSON.stringify({ user_text: userInput })
    });

    const data = await res.json();
    if (!res.ok) throw new Error(data.error || "Unknown error");

    document.getElementById("response").innerText = data.formatted_message;
  } catch (err) {
    document.getElementById("response").innerText = " " + err.message;
  }
});
