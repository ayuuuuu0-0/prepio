import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_fonts/google_fonts.dart';
import '../../core/api/api_client.dart';
import '../../core/theme/design_tokens.dart';
import '../../widgets/game/companion_hero.dart';
import '../../widgets/game/game_background.dart';
import '../../widgets/game/game_button.dart';
import '../auth/auth_provider.dart';

/// journeyProvider loads the user's journey map from the API.
final journeyProvider = FutureProvider<JourneyData>((ref) async {
  final api = ref.watch(apiClientProvider);
  final token = ref.watch(authTokenProvider);
  if (token == null || token.isEmpty) throw StateError('not authenticated');
  api.token = token;
  return api.getJourney();
});

/// JourneyScreen shows the Foundation Forest path from GET /api/v1/journey.
class JourneyScreen extends ConsumerWidget {
  const JourneyScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final journey = ref.watch(journeyProvider);

    return GameBackground(
      variant: GameBgVariant.forest,
      child: journey.when(
        loading: () => const Center(child: CircularProgressIndicator(color: PrepioColors.accent)),
        error: (e, _) => Center(child: Text('$e', style: const TextStyle(color: PrepioColors.danger))),
        data: (data) => ListView(
          padding: const EdgeInsets.all(20),
          children: [
            Text(
              '🌲 ${data.world.name}',
              textAlign: TextAlign.center,
              style: GoogleFonts.plusJakartaSans(fontSize: 24, fontWeight: FontWeight.w800, color: PrepioColors.textPrimary),
            ),
            Text(
              data.world.description,
              textAlign: TextAlign.center,
              style: GoogleFonts.nunito(color: PrepioColors.textMuted, fontWeight: FontWeight.w600),
            ),
            const SizedBox(height: 32),
            ...data.nodes.asMap().entries.map((entry) {
              final i = entry.key;
              final node = entry.value;
              final align = i.isEven ? Alignment.centerLeft : Alignment.centerRight;
              return Align(
                alignment: align,
                child: Padding(
                  padding: const EdgeInsets.symmetric(vertical: 12),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      if (node.status == 'current') ...[
                        const CompanionHero(size: 48),
                        const SizedBox(width: 8),
                      ],
                      Column(
                        children: [
                          _nodeCircle(node.status, node.nodeType),
                          const SizedBox(height: 6),
                          Text(
                            node.label,
                            style: GoogleFonts.plusJakartaSans(fontSize: 12, fontWeight: FontWeight.w700, color: PrepioColors.textBody),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              );
            }),
            const SizedBox(height: 24),
            const GameButton(label: 'Continue Prep →', onPressed: null),
          ],
        ),
      ),
    );
  }

  Widget _nodeCircle(String status, String nodeType) {
    final color = switch (status) {
      'done' => PrepioColors.success,
      'current' => PrepioColors.accent,
      _ => PrepioColors.raised,
    };
    final icon = switch (status) {
      'done' => '✓',
      'current' => '⚡',
      _ => '🔒',
    };
    final isBoss = nodeType == 'boss';
    return Container(
      width: isBoss ? 72 : 60,
      height: isBoss ? 72 : 60,
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        color: color,
        border: status == 'locked' ? Border.all(color: PrepioColors.border) : null,
      ),
      alignment: Alignment.center,
      child: Text(
        isBoss && status != 'locked' ? '👑' : icon,
        style: TextStyle(fontSize: isBoss ? 28 : 22, color: status == 'locked' ? PrepioColors.textDim : PrepioColors.bg),
      ),
    );
  }
}
