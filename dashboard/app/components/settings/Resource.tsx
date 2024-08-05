import { Text } from '@mantine/core';
import byteSize from 'byte-size';
import * as Progress from '@radix-ui/react-progress';

import { Group } from '@/components/layout/Flex';

import classes from './Resource.module.css';

interface PanelProps {
	title: string;
	usage: number;
}

interface ResourcePanelProps {
	title: string;
	used: number;
	total: number;
}

interface ResourcePanelCPUProps extends PanelProps {
	cores: number;
	threads: number;
}

const Panel = ({
	title,
	usage,
	children,
}: React.PropsWithChildren<PanelProps>) => (
	<div className={classes.panel}>
		<Group style={{ marginBottom: 16 }}>
			<Text fz={16} fw={500}>
				{title}
			</Text>
			<Text fz={22} ff="monospace">
				{usage.toFixed(2)}%
			</Text>
		</Group>
		<Progress.Root className={classes.progress} value={Math.round(usage)}>
			<Progress.Indicator
				className={classes.indicator}
				style={{ transform: `translateX(-${100 - Math.round(usage)}%)` }}
			/>
		</Progress.Root>
		{children}
	</div>
);

export const ResourcePanel = ({ title, used, total }: ResourcePanelProps) => {
	const usage = (used / total) * 100;
	return (
		<Panel title={title} usage={usage}>
			<div className={classes.group}>
				<Text>
					Used: <span>{String(byteSize(used))}</span>
				</Text>
				<Text>
					Capacity: <span>{String(byteSize(total))}</span>
				</Text>
			</div>
		</Panel>
	);
};

export const ResourcePanelCPU = ({
	title,
	usage,
	cores,
	threads,
}: ResourcePanelCPUProps) => (
	<Panel title={title} usage={usage}>
		<div className={classes.group}>
			<Text>
				Cores: <span>{String(cores)}</span>
			</Text>
			<Text>
				Threads: <span>{String(threads)}</span>
			</Text>
		</div>
	</Panel>
);
