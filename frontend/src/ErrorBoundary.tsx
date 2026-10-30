import {Component, type ErrorInfo, type ReactNode} from 'react';
import i18next from 'i18next';
import {InregistreazaEroareFrontend} from '../wailsjs/go/main/App';

interface Props {
    children: ReactNode;
    numeModul?: string;
}

interface State {
    areEroare: boolean;
    mesaj: string;
    stiva: string;
}

export class ErrorBoundary extends Component<Props, State> {
    constructor(props: Props) {
        super(props);
        this.state = {
            areEroare: false,
            mesaj: '',
            stiva: '',
        };
    }

    static getDerivedStateFromError(error: Error): State {
        return {
            areEroare: true,
            mesaj: error.message || 'unknown UI error',
            stiva: error.stack || '',
        };
    }

    componentDidCatch(error: Error, errorInfo: ErrorInfo) {
        const modul = this.props.numeModul || 'Componenta';
        const stackComplet = `${error.stack || ''}\nComponentStack:\n${errorInfo.componentStack || ''}`;
        InregistreazaEroareFrontend(modul, error.message || 'React Error', stackComplet).catch(() => {});
    }

    reincearca = () => {
        this.setState({areEroare: false, mesaj: '', stiva: ''});
    };

    render() {
        if (this.state.areEroare) {
            const numeModul = this.props.numeModul || i18next.t('eroare_ecran.modul_implicit');
            return (
                <div className="panou-eroare-componenta" role="alert">
                    <h3>
                        ⚠️ {i18next.t('eroare_ecran.titlu', {modul: numeModul})}
                    </h3>
                    <p>
                        {i18next.t('eroare_ecran.descriere')}
                    </p>
                    <div className="detalii-eroare">
                        {this.state.mesaj}
                    </div>
                    <div style={{ display: 'flex', gap: '10px' }}>
                        <button
                            type="button"
                            className="buton primar"
                            onClick={this.reincearca}
                        >
                            🔄 {i18next.t('eroare_ecran.reincearca')}
                        </button>
                        <button
                            type="button"
                            className="buton secundar"
                            onClick={() => window.location.reload()}
                        >
                            {i18next.t('eroare_ecran.reimprospateaza')}
                        </button>
                    </div>
                </div>
            );
        }

        return this.props.children;
    }
}

