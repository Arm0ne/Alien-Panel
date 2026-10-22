import{$r as e,An as t,Bi as n,Dn as r,Fi as i,Fr as a,Hi as o,Ii as s,Jr as c,Li as l,Ni as u,Pi as d,Qi as f,Un as p,Vi as m,Wn as h,Yi as g,Zi as _,Zn as v,aa as y,ai as b,ci as x,da as S,fr as C,hn as w,ii as T,ir as E,jt as D,li as O,oi as k,pa as A,qi as j,si as M,un as N,va as P,zi as F}from"./router-DhxDhlAH.js";import{d as I}from"./index-DHsLjxY8.js";function L(e){let{textColor2:t,primaryColorHover:n,primaryColorPressed:r,primaryColor:i,infoColor:a,successColor:o,warningColor:s,errorColor:l,baseColor:u,borderColor:d,opacityDisabled:f,tagColor:p,closeIconColor:m,closeIconColorHover:h,closeIconColorPressed:g,borderRadiusSmall:_,fontSizeMini:v,fontSizeTiny:y,fontSizeSmall:b,fontSizeMedium:x,heightMini:S,heightTiny:C,heightSmall:w,heightMedium:T,closeColorHover:E,closeColorPressed:D,buttonColor2Hover:O,buttonColor2Pressed:k,fontWeightStrong:A}=e;return Object.assign(Object.assign({},I),{closeBorderRadius:_,heightTiny:S,heightSmall:C,heightMedium:w,heightLarge:T,borderRadius:_,opacityDisabled:f,fontSizeTiny:v,fontSizeSmall:y,fontSizeMedium:b,fontSizeLarge:x,fontWeightStrong:A,textColorCheckable:t,textColorHoverCheckable:t,textColorPressedCheckable:t,textColorChecked:u,colorCheckable:`#0000`,colorHoverCheckable:O,colorPressedCheckable:k,colorChecked:i,colorCheckedHover:n,colorCheckedPressed:r,border:`1px solid ${d}`,textColor:t,color:p,colorBordered:`rgb(250, 250, 252)`,closeIconColor:m,closeIconColorHover:h,closeIconColorPressed:g,closeColorHover:E,closeColorPressed:D,borderPrimary:`1px solid ${c(i,{alpha:.3})}`,textColorPrimary:i,colorPrimary:c(i,{alpha:.12}),colorBorderedPrimary:c(i,{alpha:.1}),closeIconColorPrimary:i,closeIconColorHoverPrimary:i,closeIconColorPressedPrimary:i,closeColorHoverPrimary:c(i,{alpha:.12}),closeColorPressedPrimary:c(i,{alpha:.18}),borderInfo:`1px solid ${c(a,{alpha:.3})}`,textColorInfo:a,colorInfo:c(a,{alpha:.12}),colorBorderedInfo:c(a,{alpha:.1}),closeIconColorInfo:a,closeIconColorHoverInfo:a,closeIconColorPressedInfo:a,closeColorHoverInfo:c(a,{alpha:.12}),closeColorPressedInfo:c(a,{alpha:.18}),borderSuccess:`1px solid ${c(o,{alpha:.3})}`,textColorSuccess:o,colorSuccess:c(o,{alpha:.12}),colorBorderedSuccess:c(o,{alpha:.1}),closeIconColorSuccess:o,closeIconColorHoverSuccess:o,closeIconColorPressedSuccess:o,closeColorHoverSuccess:c(o,{alpha:.12}),closeColorPressedSuccess:c(o,{alpha:.18}),borderWarning:`1px solid ${c(s,{alpha:.35})}`,textColorWarning:s,colorWarning:c(s,{alpha:.15}),colorBorderedWarning:c(s,{alpha:.12}),closeIconColorWarning:s,closeIconColorHoverWarning:s,closeIconColorPressedWarning:s,closeColorHoverWarning:c(s,{alpha:.12}),closeColorPressedWarning:c(s,{alpha:.18}),borderError:`1px solid ${c(l,{alpha:.23})}`,textColorError:l,colorError:c(l,{alpha:.1}),colorBorderedError:c(l,{alpha:.08}),closeIconColorError:l,closeIconColorHoverError:l,closeIconColorPressedError:l,closeColorHoverError:c(l,{alpha:.12}),closeColorPressedError:c(l,{alpha:.18})})}var R={name:`Tag`,common:N,self:L},z={color:Object,type:{type:String,default:`default`},round:Boolean,size:String,closable:Boolean,disabled:{type:Boolean,default:void 0}},B=b(`tag`,`
 --n-close-margin: var(--n-close-margin-top) var(--n-close-margin-right) var(--n-close-margin-bottom) var(--n-close-margin-left);
 white-space: nowrap;
 position: relative;
 box-sizing: border-box;
 cursor: default;
 display: inline-flex;
 align-items: center;
 flex-wrap: nowrap;
 padding: var(--n-padding);
 border-radius: var(--n-border-radius);
 color: var(--n-text-color);
 background-color: var(--n-color);
 transition: 
 border-color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 line-height: 1;
 height: var(--n-height);
 font-size: var(--n-font-size);
`,[M(`strong`,`
 font-weight: var(--n-font-weight-strong);
 `),k(`border`,`
 pointer-events: none;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 border-radius: inherit;
 border: var(--n-border);
 transition: border-color .3s var(--n-bezier);
 `),k(`icon`,`
 display: flex;
 margin: 0 4px 0 0;
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 font-size: var(--n-avatar-size-override);
 `),k(`avatar`,`
 display: flex;
 margin: 0 6px 0 0;
 `),k(`close`,`
 margin: var(--n-close-margin);
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `),M(`round`,`
 padding: 0 calc(var(--n-height) / 3);
 border-radius: calc(var(--n-height) / 2);
 `,[k(`icon`,`
 margin: 0 4px 0 calc((var(--n-height) - 8px) / -2);
 `),k(`avatar`,`
 margin: 0 6px 0 calc((var(--n-height) - 8px) / -2);
 `),M(`closable`,`
 padding: 0 calc(var(--n-height) / 4) 0 calc(var(--n-height) / 3);
 `)]),M(`icon, avatar`,[M(`round`,`
 padding: 0 calc(var(--n-height) / 3) 0 calc(var(--n-height) / 2);
 `)]),M(`disabled`,`
 cursor: not-allowed !important;
 opacity: var(--n-opacity-disabled);
 `),M(`checkable`,`
 cursor: pointer;
 box-shadow: none;
 color: var(--n-text-color-checkable);
 background-color: var(--n-color-checkable);
 `,[x(`disabled`,[T(`&:hover`,`background-color: var(--n-color-hover-checkable);`,[x(`checked`,`color: var(--n-text-color-hover-checkable);`)]),T(`&:active`,`background-color: var(--n-color-pressed-checkable);`,[x(`checked`,`color: var(--n-text-color-pressed-checkable);`)])]),M(`checked`,`
 color: var(--n-text-color-checked);
 background-color: var(--n-color-checked);
 `,[x(`disabled`,[T(`&:hover`,`background-color: var(--n-color-checked-hover);`),T(`&:active`,`background-color: var(--n-color-checked-pressed);`)])])])]),V=Object.assign(Object.assign(Object.assign({},r.props),z),{bordered:{type:Boolean,default:void 0},checked:Boolean,checkable:Boolean,strong:Boolean,triggerClickOnClose:Boolean,onClose:[Array,Function],onMouseenter:Function,onMouseleave:Function,"onUpdate:checked":Function,onUpdateChecked:Function,internalCloseFocusable:{type:Boolean,default:!0},internalCloseIsButtonTag:{type:Boolean,default:!0},onCheckedChange:Function}),H=a(`n-tag`),U=m({name:`Tag`,props:V,slots:Object,setup(n){let i=S(null),{mergedBorderedRef:a,mergedClsPrefixRef:o,inlineThemeDisabled:s,mergedRtlRef:c,mergedComponentPropsRef:l}=h(n),d=u(()=>n.size||l?.value?.Tag?.size||`medium`),m=r(`Tag`,`-tag`,B,R,n,o);f(H,{roundRef:A(n,`round`)});function g(){if(!n.disabled&&n.checkable){let{checked:e,onCheckedChange:t,onUpdateChecked:r,"onUpdate:checked":i}=n;r&&r(!e),i&&i(!e),t&&t(!e)}}function _(e){if(n.triggerClickOnClose||e.stopPropagation(),!n.disabled){let{onClose:t}=n;t&&E(t,e)}}let v={setTextContent(e){let{value:t}=i;t&&(t.textContent=e)}},y=t(`Tag`,c,o),b=u(()=>{let{type:t,color:{color:r,textColor:i}={}}=n,o=d.value,{common:{cubicBezierEaseInOut:s},self:{padding:c,closeMargin:l,borderRadius:u,opacityDisabled:f,textColorCheckable:p,textColorHoverCheckable:h,textColorPressedCheckable:g,textColorChecked:_,colorCheckable:v,colorHoverCheckable:y,colorPressedCheckable:b,colorChecked:x,colorCheckedHover:S,colorCheckedPressed:C,closeBorderRadius:w,fontWeightStrong:T,[O(`colorBordered`,t)]:E,[O(`closeSize`,o)]:D,[O(`closeIconSize`,o)]:k,[O(`fontSize`,o)]:A,[O(`height`,o)]:j,[O(`color`,t)]:M,[O(`textColor`,t)]:N,[O(`border`,t)]:P,[O(`closeIconColor`,t)]:F,[O(`closeIconColorHover`,t)]:I,[O(`closeIconColorPressed`,t)]:L,[O(`closeColorHover`,t)]:R,[O(`closeColorPressed`,t)]:z}}=m.value,B=e(l);return{"--n-font-weight-strong":T,"--n-avatar-size-override":`calc(${j} - 8px)`,"--n-bezier":s,"--n-border-radius":u,"--n-border":P,"--n-close-icon-size":k,"--n-close-color-pressed":z,"--n-close-color-hover":R,"--n-close-border-radius":w,"--n-close-icon-color":F,"--n-close-icon-color-hover":I,"--n-close-icon-color-pressed":L,"--n-close-icon-color-disabled":F,"--n-close-margin-top":B.top,"--n-close-margin-right":B.right,"--n-close-margin-bottom":B.bottom,"--n-close-margin-left":B.left,"--n-close-size":D,"--n-color":r||(a.value?E:M),"--n-color-checkable":v,"--n-color-checked":x,"--n-color-checked-hover":S,"--n-color-checked-pressed":C,"--n-color-hover-checkable":y,"--n-color-pressed-checkable":b,"--n-font-size":A,"--n-height":j,"--n-opacity-disabled":f,"--n-padding":c,"--n-text-color":i||N,"--n-text-color-checkable":p,"--n-text-color-checked":_,"--n-text-color-hover-checkable":h,"--n-text-color-pressed-checkable":g}}),x=s?p(`tag`,u(()=>{let e=``,{type:t,color:{color:r,textColor:i}={}}=n;return e+=t[0],e+=d.value[0],r&&(e+=`a${C(r)}`),i&&(e+=`b${C(i)}`),a.value&&(e+=`c`),e}),b,n):void 0;return Object.assign(Object.assign({},v),{rtlEnabled:y,mergedClsPrefix:o,contentRef:i,mergedBordered:a,handleClick:g,handleCloseClick:_,cssVars:s?void 0:b,themeClass:x?.themeClass,onRender:x?.onRender})},render(){var e;let{mergedClsPrefix:t,rtlEnabled:n,closable:r,color:{borderColor:i}={},round:a,onRender:s,$slots:c}=this;s?.();let l=v(c.avatar,e=>e&&o(`div`,{class:`${t}-tag__avatar`},e)),u=v(c.icon,e=>e&&o(`div`,{class:`${t}-tag__icon`},e));return o(`div`,{class:[`${t}-tag`,this.themeClass,{[`${t}-tag--rtl`]:n,[`${t}-tag--strong`]:this.strong,[`${t}-tag--disabled`]:this.disabled,[`${t}-tag--checkable`]:this.checkable,[`${t}-tag--checked`]:this.checkable&&this.checked,[`${t}-tag--round`]:a,[`${t}-tag--avatar`]:l,[`${t}-tag--icon`]:u,[`${t}-tag--closable`]:r}],style:this.cssVars,onClick:this.handleClick,onMouseenter:this.onMouseenter,onMouseleave:this.onMouseleave},u||l,o(`span`,{class:`${t}-tag__content`,ref:`contentRef`},(e=this.$slots).default?.call(e)),!this.checkable&&r?o(w,{clsPrefix:t,class:`${t}-tag__close`,disabled:this.disabled,onClick:this.handleCloseClick,focusable:this.internalCloseFocusable,round:a,isButtonTag:this.internalCloseIsButtonTag,absolute:!0}):null,!this.checkable&&this.mergedBordered?o(`div`,{class:`${t}-tag__border`,style:{borderColor:i}}):null)}}),W={key:0,class:`ml-4px`},G={key:0,class:`mt-2px text-12px opacity-80`},K=m({name:`ProjectDataFreshness`,__name:`data-freshness`,props:{dataAt:{default:null},delayedAfterMinutes:{default:10},expiredAfterMinutes:{default:30},compact:{type:Boolean,default:!1}},setup(e){let t=e,r=S(Date.now()),a,o=u(()=>{if(!t.dataAt)return null;let e=new Date(t.dataAt).getTime();return Number.isFinite(e)?e:null}),c=u(()=>{if(o.value===null)return{label:`数据时间未知`,type:`warning`,detail:`中央接口没有提供有效的同步时间`};let e=r.value-o.value;if(e<-300*1e3)return{label:`数据时间未知`,type:`warning`,detail:`同步时间晚于当前时间，可能存在时钟偏差`};let n=Math.max(t.delayedAfterMinutes,1)*60*1e3,i=Math.max(t.expiredAfterMinutes,t.delayedAfterMinutes)*60*1e3;return e>i?{label:`数据过期`,type:`error`,detail:`已超过 ${Math.round(i/6e4)} 分钟未同步`}:e>n?{label:`数据延迟`,type:`warning`,detail:`已超过 ${Math.round(n/6e4)} 分钟未同步`}:{label:`数据新鲜`,type:`success`,detail:`最近一次同步在可接受范围内`}}),f=u(()=>o.value===null?``:new Date(o.value).toLocaleString(`zh-CN`,{hour12:!1}));return g(()=>{t.compact||(a=setInterval(()=>{r.value=Date.now()},6e4))}),j(()=>{a&&clearInterval(a)}),(t,r)=>{let a=U,o=D;return _(),i(o,null,{trigger:y(()=>[n(a,{size:`small`,type:c.value.type},{default:y(()=>[F(P(c.value.label)+` `,1),!e.compact&&f.value?(_(),l(`span`,W,`· `+P(f.value),1)):s(``,!0)]),_:1},8,[`type`])]),default:y(()=>[d(`div`,null,P(c.value.detail),1),f.value?(_(),l(`div`,G,`数据时间：`+P(f.value),1)):s(``,!0)]),_:1})}}});export{U as n,K as t};