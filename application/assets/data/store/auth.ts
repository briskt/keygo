export const tokenParam = 'token'
const clientIDparam = 'clientID'

export const getToken = () => getClientID() + (localStorage.getItem(tokenParam) || '')
export const getClientID = () => localStorage.getItem(clientIDparam)

init()
function init() {
    localStorage.getItem(clientIDparam) || localStorage.setItem(clientIDparam, makeRandomID())
}

const makeRandomID = () => Math.random().toString(36).slice(2)

