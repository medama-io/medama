import {
	Flex,
	FloatingIndicator,
	Group,
	Tooltip,
	UnstyledButton,
} from '@mantine/core';
import { useState } from 'react';
import { ScrollContainer } from 'react-indiana-drag-scroll';

import { IconAreaChart } from '@/components/icons/area';
import { IconBarChart } from '@/components/icons/bar';
import { InnerHeader } from '@/components/layout/InnerHeader';
import { useChartType } from '@/hooks/use-chart-type';

import { DateComboBox } from './DateSelector';
import { HeaderDataBox } from './HeaderDataBox';
import type { ChartType, StatHeaderData } from './types';

import classes from './StatsHeader.module.css';

interface StatsHeaderProps {
	stats: StatHeaderData[];
	chart: string;
}

const CHART_TYPES = [
	{
		label: 'Toggle area chart',
		value: 'area',
		icon: <IconAreaChart />,
	},
	{
		label: 'Toggle bar chart',
		value: 'bar',
		icon: <IconBarChart />,
	},
] as const;

const SegmentedChartControl = () => {
	// Segmented control for chart type
	const [rootRef, setRootRef] = useState<HTMLDivElement | null>(null);
	const [controlsRefs, setControlsRefs] = useState<
		Record<string, HTMLButtonElement | null>
	>({});

	const { setChartType, getChartType } = useChartType();
	const chartType = getChartType();

	const handleChartChange = (value: ChartType) => {
		setChartType(value);
	};

	const setControlRef = (type: ChartType) => (node: HTMLButtonElement) => {
		controlsRefs[type] = node;
		setControlsRefs(controlsRefs);
	};

	const chartTypes = CHART_TYPES.map((item) => (
		<Tooltip key={item.value} label={item.label} withArrow>
			<UnstyledButton
				key={item.value}
				className={classes.control}
				ref={setControlRef(item.value)}
				aria-label={item.label}
				onClick={() => handleChartChange(item.value)}
				data-active={chartType === item.value}
			>
				<span className={classes.controlLabel}>{item.icon}</span>
			</UnstyledButton>
		</Tooltip>
	));

	return (
		<div className={classes.toggle} ref={setRootRef}>
			<Tooltip.Group openDelay={1000}>{chartTypes}</Tooltip.Group>
			<FloatingIndicator
				component="span"
				className={classes.indicator}
				target={controlsRefs[chartType]}
				parent={rootRef}
			/>
		</div>
	);
};

const StatsHeader = ({ stats, chart }: StatsHeaderProps) => {
	return (
		<InnerHeader>
			<Flex justify="space-between" align="center" py={8}>
				<h1>Dashboard</h1>
				<DateComboBox />
			</Flex>
			<ScrollContainer>
				<Group justify="space-between" align="flex-end" mt="xs">
					<Group wrap="nowrap">
						{stats.map((stat) => (
							<HeaderDataBox
								key={stat.label}
								stat={stat}
								isActive={chart === stat.chart}
							/>
						))}
					</Group>
					<SegmentedChartControl />
				</Group>
			</ScrollContainer>
		</InnerHeader>
	);
};

export { StatsHeader };
