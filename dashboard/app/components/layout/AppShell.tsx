import { Footer } from './Footer';
import { Header } from './Header';

interface AppShellProps {
	isLoggedIn: boolean;
	children: React.ReactNode;
}

export const AppShell = ({ isLoggedIn, children }: AppShellProps) => {
	return (
		<>
			<Header isLoggedIn={isLoggedIn} />
			<main>{children}</main>
			<Footer />
		</>
	);
};
