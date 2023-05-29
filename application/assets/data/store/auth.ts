export const tokenParam = 'token'
const clientIDparam = 'clientID'

export const getToken = () => getClientID() + (localStorage.getItem(tokenParam) || '')
export const getClientID = () => localStorage.getItem(clientIDparam)

const makeRandomID = () => Math.random().toString(36).slice(2)

init()
function init() {
    localStorage.getItem(clientIDparam) || localStorage.setItem(clientIDparam, makeRandomID())
}
