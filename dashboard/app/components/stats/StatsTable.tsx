import { ActionIcon, Box, Group, Tabs, Text } from '@mantine/core';
import { Link, useNavigate, useSearchParams } from '@remix-run/react';
import {
	DataTable,
	type DataTableColumn,
	type DataTableSortStatus,
} from 'mantine-datatable';
import { useCallback, useMemo, useState } from 'react';

import { IconChevronLeft } from '@/components/icons/chevronleft';
import { IconChevronRight } from '@/components/icons/chevronright';

import { formatCount, formatDuration, formatPercentage } from './formatter';
import classes from './StatsTable.module.css';

interface DataRow {
	// Common
	visitors?: number;
	visitors_percentage?: number;
	// Mixed
	path?: string;
	bounces?: number;
	bounce_rate?: number;
	duration?: number;
	// Pages
	pageviews?: number;
	pageviews_percentage?: number;
	// Duration
	duration_upper_quartile?: number;
	duration_lower_quartile?: number;
	duration_percentage?: number;
	// Referrers
	referrer?: string;
	// Sources
	source?: string;
	// Mediums
	medium?: string;
	// Campaigns
	campaign?: string;
	// Browsers
	browser?: string;
	// Operating Systems
	os?: string;
	// Devices
	device?: string;
	// Countries
	country?: string;
	// Languages
	language?: string;
}

type QueryType = keyof typeof LABEL_MAP;

const LABEL_MAP = {
	pages: 'Pages',
	time: 'Time Spent',
	referrers: 'Referrers',
	sources: 'Sources',
	mediums: 'Mediums',
	campaigns: 'Campaigns',
	browsers: 'Browsers',
	os: 'Operating Systems',
	devices: 'Devices',
	countries: 'Countries',
	languages: 'Languages',
} as const;

const PAGE_SIZES = [10, 25, 50, 100] as const;

// Preset columns
type PresetDataKeys =
	| 'path'
	| 'visitors'
	| 'visitors_percentage'
	| 'pageviews'
	| 'pageviews_percentage'
	| 'bounce_rate'
	| 'duration';

const PRESET_COLUMNS: Record<PresetDataKeys, DataTableColumn<DataRow>> = {
	path: { accessor: 'path', title: 'Path', width: '100%' },
	visitors: { accessor: 'visitors', title: 'Visitors', sortable: true },
	visitors_percentage: {
		accessor: 'visitors_percentage',
		title: 'Visitors %',
		render: (record) => formatPercentage(record.visitors_percentage ?? 0),
	},
	pageviews: {
		accessor: 'pageviews',
		title: 'Views',
		sortable: true,
		render: (record) => formatCount(record.pageviews ?? 0),
	},
	pageviews_percentage: {
		accessor: 'pageviews_percentage',
		title: 'Views %',
		render: (record) => formatPercentage(record.pageviews_percentage ?? 0),
	},
	bounce_rate: {
		accessor: 'bounce_rate',
		title: 'Bounce %',
		render: (record) =>
			formatPercentage((record.bounces ?? 0) / (record.visitors ?? 1)),
	},
	duration: {
		accessor: 'duration',
		title: 'Duration',
		sortable: true,
		render: (record) => formatDuration(record.duration ?? 0),
	},
};

const sortBy =
	// biome-ignore lint/suspicious/noExplicitAny: Generic function.
		<T extends Record<string, any>>(key: keyof T) =>
		(a: T, b: T) =>
			a[key] > b[key] ? 1 : b[key] > a[key] ? -1 : 0;

interface QueryTableProps {
	query: QueryType;
	data: DataRow[];
}

const QueryTable = ({ query, data }: QueryTableProps) => {
	// Pagination
	const [pageSize, setPageSize] = useState<(typeof PAGE_SIZES)[number]>(10);
	const [page, setPage] = useState(1);

	// Sorting
	const [sortStatus, setSortStatus] = useState<DataTableSortStatus<DataRow>>({
		columnAccessor: 'visitors',
		direction: 'desc',
	});

	const columns = useMemo(() => getColumnsForQuery(query), [query]);

	const records = useMemo(() => {
		const from = (page - 1) * pageSize;
		const to = from + pageSize;
		const sortedData = [...data].sort(
			sortBy(sortStatus.columnAccessor as keyof DataRow),
		);
		return sortStatus.direction === 'desc'
			? sortedData.reverse().slice(from, to)
			: sortedData.slice(from, to);
	}, [data, page, pageSize, sortStatus]);

	const handlePageChange = useCallback(
		(newPage: number) => {
			const maxPage = Math.ceil(data.length / pageSize);
			setPage(Math.max(1, Math.min(newPage, maxPage)));
		},
		[data.length, pageSize],
	);

	const handlePageSizeChange = useCallback(
		(newSize: (typeof PAGE_SIZES)[number]) => {
			setPageSize(newSize);
			setPage(1);
		},
		[],
	);

	return (
		<div className={classes['table-wrapper']}>
			<div className={classes['table-header']}>
				<Text fz={14} fw={600} py={3}>
					{LABEL_MAP[query]}
				</Text>
			</div>
			<DataTable
				minHeight={300}
				highlightOnHover
				withRowBorders={false}
				records={records}
				columns={columns}
				sortStatus={sortStatus}
				onSortStatusChange={setSortStatus}
			/>
			<TablePagination
				page={page}
				pageSize={pageSize}
				totalRecords={data.length}
				onPageChange={handlePageChange}
				onPageSizeChange={handlePageSizeChange}
			/>
		</div>
	);
};

interface TablePaginationProps {
	page: number;
	pageSize: number;
	totalRecords: number;
	onPageChange: (page: number) => void;
	onPageSizeChange: (pageSize: (typeof PAGE_SIZES)[number]) => void;
}

const TablePagination = ({
	page,
	pageSize,
	totalRecords,
	onPageChange,
	onPageSizeChange,
}: TablePaginationProps) => {
	const totalPages = Math.ceil(totalRecords / pageSize);

	return (
		<Group justify="space-between" px="lg" py="sm">
			<Group>
				<span className={classes.viewspan}>View</span>
				{PAGE_SIZES.map((size) => (
					<ActionIcon
						key={size}
						variant="transparent"
						className={classes['page-size']}
						onClick={() => onPageSizeChange(size)}
						disabled={size === pageSize || totalRecords <= size}
						data-active={size === pageSize}
					>
						{size}
					</ActionIcon>
				))}
			</Group>
			<Group>
				<ActionIcon
					variant="transparent"
					className={classes['page-arrow']}
					onClick={() => onPageChange(page - 1)}
					disabled={page <= 1}
				>
					<IconChevronLeft />
				</ActionIcon>
				<span>
					Page {page} of {totalPages}
				</span>
				<ActionIcon
					variant="transparent"
					className={classes['page-arrow']}
					onClick={() => onPageChange(page + 1)}
					disabled={page >= totalPages}
				>
					<IconChevronRight />
				</ActionIcon>
			</Group>
		</Group>
	);
};

const getColumnsForQuery = (query: QueryType): DataTableColumn<DataRow>[] => {
	const commonColumns = [
		PRESET_COLUMNS.visitors,
		PRESET_COLUMNS.visitors_percentage,
		PRESET_COLUMNS.bounce_rate,
		PRESET_COLUMNS.duration,
	];

	switch (query) {
		case 'pages':
			return [
				PRESET_COLUMNS.path,
				...commonColumns,
				PRESET_COLUMNS.pageviews,
				PRESET_COLUMNS.pageviews_percentage,
			];
		case 'time':
			return [
				PRESET_COLUMNS.path,
				...commonColumns,
				{
					accessor: 'duration_lower_quartile',
					title: 'Q1 (25%)',
					sortable: true,
					render: (record) =>
						formatDuration(record.duration_lower_quartile ?? 0),
				},
				{
					accessor: 'duration_upper_quartile',
					title: 'Q3 (75%)',
					sortable: true,
					render: (record) =>
						formatDuration(record.duration_upper_quartile ?? 0),
				},
				{
					accessor: 'duration_percentage',
					title: 'Duration %',
					sortable: true,
					render: (record) =>
						formatPercentage((record.duration_percentage ?? 0) / 100),
				},
			];
		case 'referrers':
			return [
				{
					accessor: 'referrer',
					title: LABEL_MAP.referrers,
					width: '100%',
					render: (record) => record.referrer || 'Direct/None',
				},
				...commonColumns,
			];
		case 'sources':
			return [
				{
					accessor: 'source',
					title: LABEL_MAP.sources,
					width: '100%',
					render: (record) => record.source || 'Direct/None',
				},
				...commonColumns,
			];
		case 'mediums':
			return [
				{
					accessor: 'medium',
					title: LABEL_MAP.mediums,
					width: '100%',
					render: (record) => record.medium || 'Direct/None',
				},
				...commonColumns,
			];
		case 'campaigns':
			return [
				{
					accessor: 'campaign',
					title: LABEL_MAP.campaigns,
					width: '100%',
					render: (record) => record.campaign || 'Direct/None',
				},
				...commonColumns,
			];
		case 'browsers':
			return [
				{
					accessor: 'browser',
					title: LABEL_MAP.browsers,
					width: '100%',
					render: (record) => record.browser || 'Unknown',
				},
				...commonColumns,
			];
		case 'os':
			return [
				{
					accessor: 'os',
					title: LABEL_MAP.os,
					width: '100%',
					render: (record) => record.os || 'Unknown',
				},
				...commonColumns,
			];
		case 'devices':
			return [
				{
					accessor: 'device',
					title: LABEL_MAP.devices,
					width: '100%',
					render: (record) => record.device || 'Unknown',
				},
				...commonColumns,
			];
		case 'countries':
			return [
				{
					accessor: 'country',
					title: LABEL_MAP.countries,
					width: '100%',
					render: (record) => record.country || 'Unknown',
				},
				...commonColumns,
			];
		case 'languages':
			return [
				{
					accessor: 'language',
					title: LABEL_MAP.languages,
					width: '100%',
					render: (record) => record.language || 'Unknown',
				},
				...commonColumns,
			];
		default:
			throw new Error(`Invalid query: ${query}`);
	}
};

interface StatsTableProps {
	query: QueryType;
	data: DataRow[];
}

export const StatsTable = ({ query, data }: StatsTableProps) => {
	const navigate = useNavigate();
	const [searchParams] = useSearchParams();

	const handleTabChange = useCallback(
		(value: string | null) => {
			navigate(
				{
					pathname: `../${value}`,
					search: searchParams.toString(),
				},
				{ preventScrollReset: true },
			);
		},
		[navigate, searchParams],
	);

	return (
		<Tabs
			variant="unstyled"
			value={query}
			classNames={{
				root: classes['tab-root'],
				tab: classes.tab,
				list: classes['tab-list'],
			}}
			orientation="vertical"
			onChange={handleTabChange}
			keepMounted={false}
		>
			<Tabs.List>
				<div className={classes['list-wrapper']}>
					<Box
						component={Link}
						to={{
							pathname: '../',
							search: searchParams.toString(),
						}}
						className={classes.back}
					>
						<IconChevronLeft />
						<span>Go back</span>
					</Box>
					{Object.entries(LABEL_MAP).map(([key, label]) => (
						<Tabs.Tab key={key} value={key}>
							{label}
						</Tabs.Tab>
					))}
				</div>
			</Tabs.List>
			<Tabs.Panel value={query}>
				<QueryTable query={query} data={data} />
			</Tabs.Panel>
		</Tabs>
	);
};
