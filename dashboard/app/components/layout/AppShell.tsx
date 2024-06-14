import { Header } from './Header';

interface AppShellProps {
	children: React.ReactNode;
}

export const AppShell = ({ children }: AppShellProps) => {
	return (
		<>
			<Header />
			{children}
		</>
	);
};
