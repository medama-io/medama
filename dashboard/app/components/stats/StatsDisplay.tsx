import { ActionIcon, Group, Tabs, Text, UnstyledButton } from '@mantine/core';

import { IconDots } from '@/components/icons/dots';

import { countFormatter, formatDuration } from './formatter';
import classes from './StatsDisplay.module.css';

interface StatsItemProps {
	label: string;
	count: number | undefined;
	percentage: number | undefined;
	isTime?: boolean;
}

const StatsItem = ({ label, count, percentage, isTime }: StatsItemProps) => {
	const formattedValue = isTime
		? formatDuration(count ?? 0)
		: countFormatter.format(count ?? 0);
	return (
		<div className={classes['stat-item']}>
			<Group justify="space-between" pb={6}>
				<Text fz={14}>{label}</Text>
				<Text fw={600} fz={14}>
					{formattedValue}
				</Text>
			</Group>
			<div className={classes.bar} style={{ width: `${percentage}%` }} />
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
							key={item.count}
							isTime={tab.label === 'Time'}
							{...item}
						/>
					))}
				</Tabs.Panel>
			))}
			<div className={classes['button-wrapper']}>
				<UnstyledButton className={classes.button}>Load More</UnstyledButton>
			</div>
		</Tabs>
	);
};
