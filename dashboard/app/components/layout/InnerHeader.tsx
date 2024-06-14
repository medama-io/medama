import classes from './InnerHeader.module.css';

export const InnerHeader = ({ children }: { children: React.ReactNode }) => {
	return (
		<div className={classes.header}>
			<div className={classes.inner}>{children}</div>
		</div>
	);
};
