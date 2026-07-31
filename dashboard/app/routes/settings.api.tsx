import { Flex, Input } from '@mantine/core';
import { Trash2 } from 'lucide-react';
import { useCallback, useState } from 'react';
import type { ApiError } from '@/api/client';
import {
	deleteUserApiKey,
	getUserApiKey,
	regenerateUserApiKey,
} from '@/api/user';
import { Anchor } from '@/components/Anchor';
import { Button } from '@/components/Button';
import { PasswordInput } from '@/components/Input';
import { ModalChild, ModalWrapper } from '@/components/Modal';
import { SectionStack } from '@/components/settings/Section';
import type { Route } from './+types/settings.api';

type ApiSettingsModalState = 'regenerate' | 'delete' | 'hidden';

export const meta: Route.MetaFunction = () => {
	return [{ title: 'API Settings | Medama' }];
};

export const clientLoader = async () => {
	try {
		const userApiKey = await getUserApiKey();
		return {
			apiKey: userApiKey.data?.api_key,
		};
	} catch (err: unknown) {
		if ((err as ApiError)?.init?.status === 404) {
			return { apiKey: null };
		}

		throw err;
	}
};

const LearnAboutApiAnchor = () => {
	return (
		<Anchor href="https://github.com/medama-io/medama/blob/main/core/openapi.yaml">
			Learn about available API
		</Anchor>
	);
};

export default function ApiPage({ loaderData }: Route.ComponentProps) {
	const [apiKey, setApiKey] = useState(loaderData.apiKey);
	const [modalState, setModalState] = useState<ApiSettingsModalState>('hidden');

	const startRegenerateApiKey = useCallback(async () => {
		setModalState('regenerate');
	}, []);

	const confirmRegenerateApiKey = useCallback(async () => {
		const userApiKey = await regenerateUserApiKey();

		setApiKey(userApiKey.data?.api_key);
		setModalState('hidden');
	}, []);

	const startDeleteApiKey = useCallback(async () => {
		setModalState('delete');
	}, []);

	const confirmDeleteApiKey = useCallback(async () => {
		await deleteUserApiKey();

		setApiKey(null);
		setModalState('hidden');
	}, []);

	return (
		<SectionStack
			title="API Key"
			description="API key allow you to interact with Medama using your backend services."
			hasButton={false}
		>
			<div style={{ marginTop: '1rem' }}>
				{apiKey === null ? (
					<>
						<Input
							value="You haven't created an API key yet. Generate one to talk to Medama from your code."
							readOnly={true}
						/>

						<Flex mt="sm" gap="sm" justify="space-between">
							<LearnAboutApiAnchor />

							<Button onClick={confirmRegenerateApiKey}>
								Generate new API key
							</Button>
						</Flex>
					</>
				) : (
					<>
						<PasswordInput value={apiKey} readOnly={true} />
						<Flex mt="sm" gap="sm" justify="space-between" align="flex-start">
							<LearnAboutApiAnchor />
							<Flex gap="sm" justify="flex-end">
								<Button onClick={startRegenerateApiKey}>
									Generate new key
								</Button>
								<Button onClick={startDeleteApiKey} color="danger">
									<Trash2 size={16} />
								</Button>
							</Flex>
						</Flex>
					</>
				)}
			</div>

			<ModalWrapper
				opened={modalState === 'regenerate'}
				close={() => setModalState('hidden')}
			>
				<ModalChild
					title="Generate new API key"
					closeAriaLabel="Close generate new API key modal"
					description="Do you want to create new API key? Existing key will be deactivated and all services that use it will lose access to this Medama instance."
					submitLabel="Generate new API key"
					onSubmit={confirmDeleteApiKey}
					close={() => setModalState('hidden')}
				></ModalChild>
			</ModalWrapper>

			<ModalWrapper
				opened={modalState === 'delete'}
				close={() => setModalState('hidden')}
			>
				<ModalChild
					title="Delete API key"
					closeAriaLabel="Close delete API key modal"
					description="This API key will be permanently deleted. All services that use it will lose access to this Medama instance."
					submitLabel="Permanently delete API key"
					onSubmit={confirmDeleteApiKey}
					close={() => setModalState('hidden')}
					isDanger
				></ModalChild>
			</ModalWrapper>
		</SectionStack>
	);
}
