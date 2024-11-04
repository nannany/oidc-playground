window.addEventListener("message", receiveMessage, false);

function receiveMessage(e) { // e.data has client_id and session_state

    console.log('op iframe received message from RP iframe');

    var client_id = e.data.substr(0, e.data.lastIndexOf(' '));
    var session_state = e.data.substr(e.data.lastIndexOf(' ') + 1);

    console.log(client_id, session_state);

    // if message is syntactically invalid
    //     postMessage('error', e.origin) and return

    // if message comes an unexpected origin
    //     postMessage('error', e.origin) and return

    // get_op_user_agent_state() is an OP defined function
    // that returns the User Agent's login status at the OP.
    // How it is done is entirely up to the OP.
    var opuas = get_op_user_agent_state();

    // Here, the session_state is calculated in this particular way,
    // but it is entirely up to the OP how to do it under the
    // requirements defined in this specification.
    var ss = CryptoJS.SHA256(client_id + ' ' + e.origin + ' ' +
        opuas + ' ' + salt) + "." + salt;

    var stat = '';
    if (session_state === ss) {
        stat = 'unchanged';
    } else {
        stat = 'changed';
    }

    e.source.postMessage(stat, e.origin);
};
