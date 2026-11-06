import {useTranslation} from 'react-i18next';
export function SelectieRezumat({persoane,deLa,panaLa}:{persoane:{id:string;marca:string}[];deLa?:string;panaLa?:string}){const {t}=useTranslation();return <p className="ux30-selectie" role="status">{t('ux30.selectie',{numar:persoane.length,deLa:deLa||'--',panaLa:panaLa||'--',marci:persoane.map(p=>p.marca).join(', ')||'--'})}</p>;}
