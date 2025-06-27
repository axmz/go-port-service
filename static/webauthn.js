document.getElementById('registerBtn').onclick = async function () {
    const email = document.getElementById('email').value;
    if (!email) return output('Please enter a valid email.');

    // Begin registration (fetch challenge/options from backend)
    let res = await fetch('/api/webauth/register/begin', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email })
    });
    if (!res.ok) return output('Failed to start registration');
    let data = await res.json();
    let options = data.data;

    // Convert challenge and user.id to Uint8Array
    options.publicKey.challenge = base64urlToBuffer(options.publicKey.challenge);
    options.publicKey.user.id = base64urlToBuffer(options.publicKey.user.id);

    // Call WebAuthn API
    let cred;
    try {
        cred = await navigator.credentials.create(options);
    } catch (e) {
        return output('Registration failed: ' + e);
    }

    // Send attestation to backend
    let attestation = {
        id: cred.id,
        rawId: bufferToBase64url(cred.rawId),
        type: cred.type,
        response: {
            clientDataJSON: bufferToBase64url(cred.response.clientDataJSON),
            attestationObject: bufferToBase64url(cred.response.attestationObject)
        }
    };

    res = await fetch('/api/webauth/register/finish', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(attestation)
    });
    if (!res.ok) return output('Registration finish failed');

    // Show output and countdown before redirect
    let seconds = 3;
    output(`Registration successful! Redirecting to /private in ${seconds} seconds...`);
    const countdown = setInterval(() => {
        seconds--;
        if (seconds > 0) {
            output(`Registration successful! Redirecting to /private in ${seconds} seconds...`);
        } else {
            clearInterval(countdown);
            window.location.href = "/private";
        }
    }, 1000);
};

document.getElementById('loginBtn').onclick = async function () {
    const email = document.getElementById('email').value;
    if (!email) return output('Please enter a valid email.');

    // Begin login (fetch challenge/options from backend)
    let res = await fetch('/api/webauth/login/begin', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email })
    });
    if (!res.ok) return output('Failed to start login');
    let data = await res.json();
    let options = data.data;

    // Convert challenge and allowCredentials.id to Uint8Array
    options.publicKey.challenge = base64urlToBuffer(options.publicKey.challenge);
    if (options.publicKey.allowCredentials) {
        options.publicKey.allowCredentials = options.publicKey.allowCredentials.map(cred => ({
            ...cred,
            id: base64urlToBuffer(cred.id)
        }));
    }

    // Call WebAuthn API
    let assertion;
    try {
        assertion = await navigator.credentials.get(options);
    } catch (e) {
        return output('Login failed: ' + e);
    }

    // Send assertion to backend
    let authData = {
        id: assertion.id,
        rawId: bufferToBase64url(assertion.rawId),
        type: assertion.type,
        response: {
            clientDataJSON: bufferToBase64url(assertion.response.clientDataJSON),
            authenticatorData: bufferToBase64url(assertion.response.authenticatorData),
            signature: bufferToBase64url(assertion.response.signature),
            userHandle: assertion.response.userHandle ? bufferToBase64url(assertion.response.userHandle) : null
        }
    };

    res = await fetch('/api/webauth/login/finish', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(authData)
    });
    if (!res.ok) return output('Login finish failed');

    // Show output and countdown before redirect
    let seconds = 3;
    output(`Login successful! Redirecting to /private in ${seconds} seconds...`);
    const countdown = setInterval(() => {
        seconds--;
        if (seconds > 0) {
            output(`Login successful! Redirecting to /private in ${seconds} seconds...`);
        } else {
            clearInterval(countdown);
            window.location.href = "/private";
        }
    }, 1000);
};

// Helper functions
function output(msg) {
    document.getElementById('output').textContent = msg;
}
function base64urlToBuffer(base64url) {
    // Pad base64 string and convert to Uint8Array
    let base64 = base64url.replace(/-/g, '+').replace(/_/g, '/');
    while (base64.length % 4) base64 += '=';
    let str = atob(base64);
    let bytes = new Uint8Array(str.length);
    for (let i = 0; i < str.length; ++i) bytes[i] = str.charCodeAt(i);
    return bytes.buffer;
}
function bufferToBase64url(buffer) {
    let bytes = new Uint8Array(buffer);
    let str = '';
    for (let i = 0; i < bytes.byteLength; ++i) str += String.fromCharCode(bytes[i]);
    return btoa(str).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}