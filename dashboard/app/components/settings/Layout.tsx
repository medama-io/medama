import classes from './Layout.module.css';
import { Sidebar } from './Sidebar';

interface SettingsLayoutProps {
	children: React.ReactNode;
}

export const SettingsLayout = ({ children }: SettingsLayoutProps) => {
	return (
		<div className={classes.wrapper}>
			<Sidebar />
			<div className={classes.flex}>{children}</div>
		</div>
	);
};
