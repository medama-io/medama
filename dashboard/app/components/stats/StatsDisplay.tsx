import { Tabs, UnstyledButton } from '@mantine/core';
import { Link, useSearchParams } from '@remix-run/react';

import classes from './StatsDisplay.module.css';
import { StatsItem } from './StatsItem';

interface StatsValue {
	label: string;
	count?: number;
	percentage?: number;
}

export interface StatsTab {
	label: string;
	items: StatsValue[];
}

interface StatsDisplayProps {
	data: StatsTab[];
}

export const StatsDisplay = ({ data }: StatsDisplayProps) => {
	const [searchParams] = useSearchParams();
	const defaultValue = data[0]?.label ?? '';

	return (
		<Tabs
			variant="unstyled"
			defaultValue={defaultValue}
			classNames={{
				root: classes.root,
				tab: classes.tab,
				list: classes.list,
			}}
		>
			<Tabs.List>
				{data.map((tab) => (
					<Tabs.Tab key={tab.label} value={tab.label} aria-label={tab.label}>
						{tab.label}
					</Tabs.Tab>
				))}
			</Tabs.List>

			{data.map((tab) => (
				<Tabs.Panel key={tab.label} value={tab.label}>
					<div style={{ minHeight: 306 }}>
						{tab.items.map((item) => (
							<StatsItem key={item.label} tab={tab.label} {...item} />
						))}
					</div>
					<div className={classes.more}>
						<UnstyledButton
							component={Link}
							to={{
								pathname: `./${tab.label.toLowerCase()}`,
								search: `?${searchParams.toString()}`,
							}}
							prefetch="intent"
							preventScrollReset
							className={classes.button}
							aria-label={`Load more ${tab.label} stats.`}
						>
							Load More
						</UnstyledButton>
					</div>
				</Tabs.Panel>
			))}
		</Tabs>
	);
};
