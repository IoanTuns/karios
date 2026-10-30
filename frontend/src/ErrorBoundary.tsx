import {Component, type ErrorInfo, type ReactNode} from 'react';
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
            return (
                <div className="panou-eroare-componenta" role="alert">
                    <h3>
                        ⚠️ A apărut o problemă la afișarea acestui ecran ({this.props.numeModul || 'Modul'})
                    </h3>
                    <p>
                        Eroarea a fost captată și înregistrată automat în jurnalul local de diagnoză. Baza ta de date și datele de pontaj sunt în siguranță.
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
                            🔄 Reîncearcă afișarea
                        </button>
                        <button
                            type="button"
                            className="buton secundar"
                            onClick={() => window.location.reload()}
                        >
                            Reîmprospătează aplicația
                        </button>
                    </div>
                </div>
            );
        }

        return this.props.children;
    }
}

