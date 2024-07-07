import { Tabs, Text, UnstyledButton } from '@mantine/core';
import { Link, useSearchParams } from '@remix-run/react';

import classes from './StatsDisplay.module.css';
import { StatsItem } from './StatsItem';
import type { StatsTab } from './types';

interface StatsDisplayProps {
	data: StatsTab[];
}

interface LoadMoreButtonProps {
	tab: StatsTab;
	searchParams: URLSearchParams;
}

export const StatsDisplay = ({ data }: StatsDisplayProps) => {
	const [searchParams] = useSearchParams();

	return (
		<Tabs
			variant="unstyled"
			defaultValue={data[0]?.label}
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
					<div className={classes.items} data-empty={tab.items.length === 0}>
						{tab.items.map((item) => (
							<StatsItem key={item.label} tab={tab.label} {...item} />
						))}
						{tab.items.length === 0 && (
							<Text c="gray" fz={14}>
								No records found...
							</Text>
						)}
					</div>
					<LoadMoreButton tab={tab} searchParams={searchParams} />
				</Tabs.Panel>
			))}
		</Tabs>
	);
};

const LoadMoreButton = ({ tab, searchParams }: LoadMoreButtonProps) => (
	<div className={classes.more}>
		<UnstyledButton
			component={Link}
			to={{
				pathname: `./${tab.label.toLowerCase()}`,
				search: searchParams.toString(),
			}}
			prefetch="intent"
			preventScrollReset
			className={classes.button}
			aria-label={`Load more ${tab.label} stats.`}
		>
			Load More
		</UnstyledButton>
	</div>
);
