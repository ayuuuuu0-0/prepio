import 'package:flutter/material.dart';
import '../../core/theme/design_tokens.dart';

/// CompanionReaction drives correct/wrong/idle companion states.
enum CompanionReaction { idle, correct, wrong }

/// CompanionHero shows a large animated companion avatar with glow ring.
class CompanionHero extends StatefulWidget {
  const CompanionHero({
    super.key,
    this.name,
    this.species,
    this.size = 96,
    this.reaction = CompanionReaction.idle,
  });

  final String? name;
  final String? species;
  final double size;
  final CompanionReaction reaction;

  @override
  State<CompanionHero> createState() => _CompanionHeroState();
}

class _CompanionHeroState extends State<CompanionHero> with SingleTickerProviderStateMixin {
  late final AnimationController _controller;
  late final Animation<double> _bounce;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: widget.reaction == CompanionReaction.correct
          ? const Duration(milliseconds: 700)
          : const Duration(milliseconds: 3000),
    );
    if (widget.reaction == CompanionReaction.idle) {
      _controller.repeat(reverse: true);
    } else {
      _controller.forward();
    }
    final bounceEnd = widget.reaction == CompanionReaction.correct ? -14.0 : -6.0;
    _bounce = Tween<double>(begin: 0, end: bounceEnd)
        .animate(CurvedAnimation(parent: _controller, curve: Curves.easeInOut));
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final visual = companionFor(name: widget.name, species: widget.species);

    return AnimatedBuilder(
      animation: _bounce,
      builder: (context, child) => Transform.translate(offset: Offset(0, _bounce.value), child: child),
      child: Stack(
        alignment: Alignment.center,
        children: [
          Container(
            width: widget.size * 1.2,
            height: widget.size * 1.2,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              boxShadow: [BoxShadow(color: visual.glow.withValues(alpha: 0.4), blurRadius: 24)],
            ),
          ),
          Container(
            width: widget.size,
            height: widget.size,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              gradient: LinearGradient(
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
                colors: [visual.glow.withValues(alpha: 0.2), PrepioColors.surface],
              ),
              border: Border.all(color: Colors.white.withValues(alpha: 0.1)),
              boxShadow: [
                BoxShadow(color: visual.glow.withValues(alpha: 0.25), blurRadius: 16),
              ],
            ),
            alignment: Alignment.center,
            child: Text(visual.emoji, style: TextStyle(fontSize: widget.size * 0.42)),
          ),
        ],
      ),
    );
  }
}
