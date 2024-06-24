import { Group, Tabs, Text, UnstyledButton } from '@mantine/core';
import { Link, useSearchParams } from '@remix-run/react';

import { formatCount, formatDuration } from './formatter';
import classes from './StatsDisplay.module.css';

interface StatsItem {
	label: string;
	count: number | undefined;
	percentage: number | undefined;
}

interface StatsItemProps extends StatsItem {
	tab: string;
}

export interface StatsTab {
	label: string;
	items: StatsItem[];
}

interface StatsDisplayProps {
	data: StatsTab[];
}

const StatsItem = ({ label, count, percentage, tab }: StatsItemProps) => {
	const [searchParams, setSearchParams] = useSearchParams();

	const formattedValue =
		tab === 'Time' ? formatDuration(count ?? 0) : formatCount(count ?? 0);

	const handleFilter = () => {
		if (tab !== 'Time') {
			const params = new URLSearchParams(searchParams);

			const filterMap: Record<string, string> = {
				Referrers: 'referrer',
				Sources: 'utm_source',
				Mediums: 'utm_medium',
				Campaigns: 'utm_campaign',
				Browsers: 'browser',
				OS: 'os',
				Devices: 'device',
				Countries: 'country',
				Languages: 'language',
			};
			const filter = filterMap[tab] ?? 'path';

			params.append(`${filter}[eq]`, label);
			setSearchParams(params, {
				preventScrollReset: true,
			});
		}
	};

	return (
		<UnstyledButton className={classes['stat-item']} onClick={handleFilter}>
			<Group justify="space-between" pb={6}>
				<Text fz={14}>{label}</Text>
				<Text fw={600} fz={14}>
					{formattedValue}
				</Text>
			</Group>
			<div
				className={classes.bar}
				style={{ width: `${(percentage ?? 0) * 100}%` }}
			/>
		</UnstyledButton>
	);
};

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
					<Tabs.Tab key={tab.label} value={tab.label}>
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
						>
							Load More
						</UnstyledButton>
					</div>
				</Tabs.Panel>
			))}
		</Tabs>
	);
};
