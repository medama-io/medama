import classes from './Anchor.module.css';

interface AnchorProps extends React.AnchorHTMLAttributes<HTMLAnchorElement> {
	children: React.ReactNode;
}

const Anchor = ({ children, className, ...rest }: AnchorProps) => {
	const isExternal = rest.href?.toString().startsWith('http');
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
