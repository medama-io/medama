import { Link } from '@remix-run/react';

import classes from './WebsiteCard.module.css';

interface WebsiteCardProps {
	website: {
		hostname: string;
		summary?: {
			visitors: number;
		};
	};
}

export const WebsiteCard = ({ website }: WebsiteCardProps) => {
	return (
		<Link
			to={`/${website.hostname}`}
			prefetch="intent"
			className={classes.card}
			tabIndex={0}
			aria-label={`Visit ${website.hostname}`}
			aria-describedby={`${website.hostname}-visitors`}
		>
			<p>{website.hostname}</p>
			<span id={`${website.hostname}-visitors`}>
				{website.summary?.visitors ?? 'N/A'} visitors in the last 24 hours
			</span>
		</Link>
	);
};
