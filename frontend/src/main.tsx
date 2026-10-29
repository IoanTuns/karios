import {StrictMode} from 'react';
import {createRoot} from 'react-dom/client';
import './style.css';
import App from './App';
import {PoartaCont} from './Conturi';
import {initI18n} from './i18n';
import {InregistreazaEroareFrontend} from '../wailsjs/go/main/App';

// Captare globală a excepțiilor JavaScript neprinse (D043, F16)
window.addEventListener('error', (event) => {
    const mesaj = event.message || 'uncaught script error';
    const sursa = `${event.filename || 'unknown'}:${event.lineno || 0}`;
    const stack = event.error?.stack || '';
    InregistreazaEroareFrontend(sursa, mesaj, stack).catch(() => {});
});

window.addEventListener('unhandledrejection', (event) => {
    const motiv = event.reason instanceof Error ? event.reason.message : String(event.reason);
    const stack = event.reason instanceof Error ? event.reason.stack || '' : '';
    InregistreazaEroareFrontend('PromiseRespins', motiv, stack).catch(() => {});
});

initI18n().then(() => {
    createRoot(document.getElementById('root')!).render(
        <StrictMode>
            <PoartaCont><App/></PoartaCont>
        </StrictMode>
    );
});
