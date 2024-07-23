import {
	Flex,
	FloatingIndicator,
	Group,
	Tooltip,
	UnstyledButton,
} from '@mantine/core';
import { useMediaQuery } from '@mantine/hooks';
import { useParams } from '@remix-run/react';
import { CalendarFold } from 'lucide-react';
import type React from 'react';
import { useState } from 'react';
import { ScrollContainer } from 'react-indiana-drag-scroll';

import { datePickerClasses, DatePickerRange } from '@/components/DatePicker';
import { DropdownSelect } from '@/components/DropdownSelect';
import { IconAreaChart } from '@/components/icons/area';
import { IconBarChart } from '@/components/icons/bar';
import { InnerHeader } from '@/components/layout/InnerHeader';
import { useChartType } from '@/hooks/use-chart-type';
import { useDisclosure } from '@/hooks/use-disclosure';

import { HeaderDataBox } from './HeaderDataBox';
import type { ChartType, StatHeaderData } from './types';

import classes from './StatsHeader.module.css';

interface StatsHeaderProps {
	stats: StatHeaderData[];
	chart: string;
	websites: string[];
}

const DATE_GROUP_END_VALUES = ['yesterday', '30d', 'year', 'all'];

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

const StatsHeader = ({ stats, chart, websites }: StatsHeaderProps) => {
	const { hostname } = useParams();
	const [dateOpened, { toggle: toggleDate }] = useDisclosure(false);
	const hideWebsiteSelector = useMediaQuery('(36em < width < 62em)');

	// Convert websites array to object with same key-val for DropdownSelect
	const websitesRecord = Object.fromEntries(
		websites.map((website) => [website, website]),
	);

	const DATE_PRESETS: Record<string, React.ReactNode> = {
		today: 'Today',
		yesterday: 'Yesterday',
		'12h': 'Previous 12 hours',
		'24h': 'Previous 24 hours',
		'72h': 'Previous 72 hours',
		'7d': 'Previous 7 days',
		'14d': 'Previous 14 days',
		'30d': 'Previous 30 days',
		quarter: 'Previous quarter',
		half: 'Previous half year',
		year: 'Previous year',
		all: 'All time',
		custom: (
			<button
				type="button"
				onClick={toggleDate}
				className={datePickerClasses.button}
			>
				Custom range
			</button>
		),
	};

	return (
		<InnerHeader>
			<Flex className={classes.title}>
				<h1>Dashboard</h1>
				<Group className={classes.dropdowns}>
					{!hideWebsiteSelector && (
						<DropdownSelect
							records={websitesRecord}
							defaultValue={hostname ?? ''}
							defaultLabel="Unknown"
							ariaLabel="Select website"
							shouldUseNavigate
						/>
					)}
					<DropdownSelect
						records={DATE_PRESETS}
						defaultValue="today"
						defaultLabel="Custom range"
						ariaLabel="Select date range"
						icon={CalendarFold}
						searchParamKey="period"
						separatorValues={DATE_GROUP_END_VALUES}
					/>
					<DatePickerRange open={dateOpened} setOpen={toggleDate} />
				</Group>
			</Flex>
			<ScrollContainer>
				<Group justify="space-between" align="flex-end" mt="xs">
					<Group wrap="nowrap" p={4}>
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
