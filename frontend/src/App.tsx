import {useEffect, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {Informatii} from '../wailsjs/go/main/App';
import {main} from '../wailsjs/go/models';

function App() {
    const {t} = useTranslation();
    const [info, setInfo] = useState<main.Informatii>();

    useEffect(() => {
        Informatii().then(setInfo);
    }, []);

    return (
        <main className="ecran">
            <h1>{t('aplicatie.titlu')}</h1>
            <p>{t('aplicatie.subtitlu')}</p>
            {info && <p className="versiune">{t('aplicatie.versiune', {versiune: info.versiune})}</p>}
        </main>
    );
}

export default App;
