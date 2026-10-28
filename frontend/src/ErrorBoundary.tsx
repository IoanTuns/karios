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
            mesaj: error.message || 'Eroare necunoscută în interfață',
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
                <div className="panou-eroare-componenta" role="alert" style={{
                    padding: '28px',
                    margin: '20px',
                    background: '#fff6f5',
                    border: '1px solid #f8d7da',
                    borderRadius: '10px',
                    color: '#721c24'
                }}>
                    <h3 style={{ margin: '0 0 10px', display: 'flex', alignItems: 'center', gap: '8px' }}>
                        ⚠️ A apărut o problemă la afișarea acestui ecran ({this.props.numeModul || 'Modul'})
                    </h3>
                    <p style={{ margin: '0 0 14px', fontSize: '13px', color: '#555' }}>
                        Eroarea a fost captată și înregistrată automat în jurnalul local de diagnoză. Baza ta de date și datele de pontaj sunt în siguranță.
                    </p>
                    <div style={{
                        background: '#ffffff',
                        border: '1px solid #ebd0d3',
                        padding: '10px 14px',
                        borderRadius: '6px',
                        fontFamily: 'monospace',
                        fontSize: '12px',
                        marginBottom: '16px',
                        overflowX: 'auto',
                        whiteSpace: 'pre-wrap'
                    }}>
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

