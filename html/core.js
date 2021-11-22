var net = {
    server: "http://localhost/router?",

    post: async function (handle = '', data) {
        var obj = {};
        if (data instanceof FormData) {
            data.forEach((value, key) => obj[key] = value);
        } else {
            if (!!data) {
                obj = data
            }
        }

        obj['handle'] = handle
        const response = await fetch(this.server, {
            method: 'POST',
            mode: 'cors',
            cache: 'no-cache',
            credentials: 'same-origin',
            headers: {
                'Content-Type': 'application/json; charset=UTF-8'
            },
            redirect: 'follow',
            referrerPolicy: 'no-referrer',
            body: JSON.stringify(obj)
        });
        return response.json();
    }
}
