import { is } from 'valibot';
import classes from './Anchor.module.css';

interface AnchorProps extends React.AnchorHTMLAttributes<HTMLAnchorElement> {
	children: React.ReactNode;
	isExternal?: boolean;
}

const Anchor = ({ children, className, isExternal, ...rest }: AnchorProps) => {
	return (
		<a
			className={className ? className : classes.anchor}
			target={isExternal ? '_blank' : undefined}
			rel={isExternal ? 'noopener noreferrer' : undefined}
			{...rest}
		>
			{children}
		</a>
	);
};

export { Anchor };
