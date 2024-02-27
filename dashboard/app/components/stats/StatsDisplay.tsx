import { ActionIcon, Group, Tabs, Text, UnstyledButton } from '@mantine/core';
import { Link, useSearchParams } from '@remix-run/react';

import { IconDots } from '@/components/icons/dots';

import { formatCount, formatDuration } from './formatter';
import classes from './StatsDisplay.module.css';

interface StatsItemProps {
	label: string;
	count: number | undefined;
	percentage: number | undefined;
	isTime?: boolean;
}

const StatsItem = ({ label, count, percentage, isTime }: StatsItemProps) => {
	const formattedValue = isTime ? formatDuration(count) : formatCount(count);
	return (
		<div className={classes['stat-item']}>
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
		</div>
	);
};

export interface StatsTab {
	label: string;
	items: StatsItemProps[];
}

interface StatsDisplayProps {
	data: StatsTab[];
}

export const StatsDisplay = ({ data }: StatsDisplayProps) => {
	const [searchParams] = useSearchParams();
	return (
		<Tabs
			variant="unstyled"
			defaultValue={data[0]?.label}
			classNames={{
				root: classes['tab-root'],
				tab: classes.tab,
			}}
		>
			<Group justify="space-between" className={classes['tab-list']}>
				<Tabs.List>
					{data.map((tab) => (
						<Tabs.Tab key={tab.label} value={tab.label}>
							{tab.label}
						</Tabs.Tab>
					))}
				</Tabs.List>
				<ActionIcon variant="transparent">
					<IconDots />
				</ActionIcon>
			</Group>

			{data.map((tab) => (
				<Tabs.Panel key={tab.label} value={tab.label}>
					{tab.items.map((item) => (
						<StatsItem
							key={item.label}
							isTime={tab.label === 'Time'}
							{...item}
						/>
					))}
					<div className={classes['button-wrapper']}>
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
